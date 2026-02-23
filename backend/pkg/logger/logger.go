package logger

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	myctx "github.com/devdavidalonso/cecor/backend/pkg/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithContext adds the logger to the provided context
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, myctx.LoggerKey, logger)
}

// FromContext extracts the logger from the context
// Returns the logger and a boolean indicating if the logger was found
func FromContext(ctx context.Context) (Logger, bool) {
	logger, ok := ctx.Value(myctx.LoggerKey).(Logger)
	return logger, ok
}

// GetLoggerOrDefault gets the logger from context or returns a default logger if not found
func GetLoggerOrDefault(ctx context.Context) Logger {
	logger, ok := FromContext(ctx)
	if !ok {
		return NewLogger() // Fallback to a new logger
	}
	return logger
}

// configureLogstash adiciona um novo sink para Logstash
// Em backend/pkg/logger/logger.go
func configureLogstash(writers *[]zapcore.WriteSyncer) {
	if os.Getenv("ENABLE_ELK") == "true" {
		logstashAddr := os.Getenv("LOGSTASH_ADDR")
		if logstashAddr == "" {
			logstashAddr = "logstash:5000"
		}

		// Usar um pool de conexões ou reconexão automática
		conn, err := net.Dial("tcp", logstashAddr)
		if err != nil {
			fmt.Printf("Failed to connect to Logstash: %v\n", err)
			return
		}

		// Adicionar um wrapper que tente reconectar em caso de falha
		writer := &reconnectingWriter{
			addr: logstashAddr,
			conn: conn,
		}
		*writers = append(*writers, zapcore.AddSync(writer))
	}
}

func (w *reconnectingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn == nil {
		w.conn, err = net.Dial("tcp", w.addr)
		if err != nil {
			return 0, err
		}
	}

	n, err = w.conn.Write(p)
	if err != nil {
		w.conn.Close()
		w.conn = nil
	}
	return
}

func (w *reconnectingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}

// Implementar um writer com reconexão automática
type reconnectingWriter struct {
	addr string
	conn net.Conn
	mu   sync.Mutex
}

// Logger é uma interface para logging
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Sync() error
}

// func (l Logger) NewLogger() Logger {
// 	panic("unimplemented")
// }

// zapLogger implementa a interface Logger usando zap
type zapLogger struct {
	logger *zap.SugaredLogger
}

// Em backend/pkg/logger/logger.go
func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	// Adicionar campos padrão para todos os logs
	standardFields := []interface{}{
		"service", "lar-cepprs-backend",
		"version", os.Getenv("APP_VERSION"),
		"env", os.Getenv("APP_ENV"),
	}

	// Combinar campos padrão com os fornecidos
	allFields := append(standardFields, keysAndValues...)

	l.logger.Infow(msg, allFields...)
}

// NewLogger cria uma nova instância de Logger
func NewLogger() Logger {
	// Determinar se estamos em produção
	isProduction := os.Getenv("APP_ENV") == "production"

	// Configurar encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Escolher encoder e nível de log com base no ambiente
	var encoder zapcore.Encoder
	var logLevel zapcore.Level

	if isProduction {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		logLevel = zapcore.InfoLevel
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		logLevel = zapcore.DebugLevel
	}

	// Preparar os writers para os logs
	writers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}

	// Adicionar Logstash como um writer adicional, se configurado
	configureLogstash(&writers)

	// Criar core com múltiplos escritores
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		logLevel,
	)

	// Criar logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &zapLogger{
		logger: logger.Sugar(),
	}
}

// Debug loga em nível de debug
func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// // Info loga em nível de info
// func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
// 	l.logger.Infow(msg, keysAndValues...)
// }

// Warn loga em nível de aviso
func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

// Error loga em nível de erro
func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Fatal loga em nível fatal e encerra a aplicação
func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

// Sync finaliza o logger corretamente
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
