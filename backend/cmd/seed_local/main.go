package main

import (
	"fmt"
	"log"

	"github.com/devdavidalonso/cecor/backend/internal/config"
	"github.com/devdavidalonso/cecor/backend/internal/repository/postgres"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err := postgres.MigrateDB(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("cecor2024!"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	if err := runSeed(db, string(passwordHash)); err != nil {
		log.Fatalf("failed to seed data: %v", err)
	}

	log.Println("seed local concluido com sucesso")
}

func runSeed(db *gorm.DB, passwordHash string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// No schema legado local existe dependência circular entre users.profile_id <-> user_profiles.user_id.
		// Para seed idempotente, removemos e recriamos constraints dentro da transação.
		if err := tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_profile`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`ALTER TABLE user_profiles DROP CONSTRAINT IF EXISTS user_profiles_user_id_fkey`).Error; err != nil {
			return err
		}

		// 1) Upsert usuários necessários para telas de admin/professor/aluno.
		if err := upsertUser(tx, "Admin CECOR", "admin.cecor@cecor.org", "11111111101", "(11) 91111-1111", "administrator", 1, passwordHash); err != nil {
			return err
		}
		if err := upsertUser(tx, "Carlos Ferreira", "carlos.ferreira@cecor.org", "23456789100", "(11) 95432-1098", "teacher", 2, passwordHash); err != nil {
			return err
		}
		if err := upsertUser(tx, "Mariana Costa", "mariana.costa@cecor.org", "34567891200", "(11) 94321-0987", "professor", 2, passwordHash); err != nil {
			return err
		}
		if err := upsertUser(tx, "Pedro Almeida", "pedro.almeida@cecor.org", "56789012300", "(11) 93210-9876", "student", 3, passwordHash); err != nil {
			return err
		}
		if err := upsertUser(tx, "Ana Clara Souza", "ana.clara@cecor.org", "67890123400", "(11) 92109-8765", "student", 3, passwordHash); err != nil {
			return err
		}

		adminID, err := getUserIDByEmail(tx, "admin.cecor@cecor.org")
		if err != nil {
			return err
		}
		profID, err := getUserIDByEmail(tx, "carlos.ferreira@cecor.org")
		if err != nil {
			return err
		}
		studentID, err := getUserIDByEmail(tx, "pedro.almeida@cecor.org")
		if err != nil {
			return err
		}

		// 2) Garante catálogo de perfis em IDs fixos esperados pelo sistema atual.
		if err := upsertProfile(tx, 1, adminID, "administrator", "admin", "Administrador"); err != nil {
			return err
		}
		if err := upsertProfile(tx, 2, profID, "teacher", "teacher", "Professor"); err != nil {
			return err
		}
		if err := upsertProfile(tx, 3, studentID, "student", "student", "Aluno"); err != nil {
			return err
		}

		// 3) Reforça profile_id correto nos usuários seedados.
		if err := tx.Exec(`
			UPDATE users
			SET profile_id = CASE
				WHEN LOWER(email) IN ('admin.cecor@cecor.org') THEN 1
				WHEN LOWER(email) IN ('carlos.ferreira@cecor.org','mariana.costa@cecor.org') THEN 2
				WHEN LOWER(email) IN ('pedro.almeida@cecor.org','ana.clara@cecor.org') THEN 3
				ELSE profile_id
			END
			WHERE LOWER(email) IN (
				'admin.cecor@cecor.org',
				'carlos.ferreira@cecor.org',
				'mariana.costa@cecor.org',
				'pedro.almeida@cecor.org',
				'ana.clara@cecor.org'
			)
		`).Error; err != nil {
			return err
		}

		// 4) Students vinculados aos usuários com profile_id=3.
		if err := upsertStudent(tx, "pedro.almeida@cecor.org", "20260001"); err != nil {
			return err
		}
		if err := upsertStudent(tx, "ana.clara@cecor.org", "20260002"); err != nil {
			return err
		}

		// 5) Recria constraints removidas no início.
		if err := tx.Exec(`
			ALTER TABLE user_profiles
			ADD CONSTRAINT user_profiles_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES users(id)
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE users
			ADD CONSTRAINT fk_users_profile
			FOREIGN KEY (profile_id) REFERENCES user_profiles(id)
		`).Error; err != nil {
			return err
		}

		return nil
	})
}

func upsertUser(tx *gorm.DB, name, email, cpf, phone, profileText string, profileID int, passwordHash string) error {
	var count int64
	if err := tx.Raw(`SELECT COUNT(1) FROM users WHERE LOWER(email)=LOWER(?)`, email).Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return tx.Exec(`
			UPDATE users
			SET name = ?,
				cpf = ?,
				phone = ?,
				profile = ?,
				profile_id = ?,
				password = ?,
				active = true,
				updated_at = NOW()
			WHERE LOWER(email)=LOWER(?)
		`, name, cpf, phone, profileText, profileID, passwordHash, email).Error
	}

	return tx.Exec(`
		INSERT INTO users (
			name, email, password, profile, cpf, phone, active, profile_id, birth_date, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, true, ?, NOW() - INTERVAL '20 years', NOW(), NOW())
	`, name, email, passwordHash, profileText, cpf, phone, profileID).Error
}

func getUserIDByEmail(tx *gorm.DB, email string) (int64, error) {
	var id int64
	if err := tx.Raw(`SELECT id FROM users WHERE LOWER(email)=LOWER(?) LIMIT 1`, email).Scan(&id).Error; err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("user not found for email: %s", email)
	}
	return id, nil
}

func upsertProfile(tx *gorm.DB, id int, userID int64, name, profileType, description string) error {
	var count int64
	if err := tx.Raw(`SELECT COUNT(1) FROM user_profiles WHERE id = ?`, id).Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return tx.Exec(`
			UPDATE user_profiles
			SET user_id = ?,
				name = ?,
				profile_type = ?,
				description = ?,
				is_primary = true,
				is_active = true,
				start_date = COALESCE(start_date, CURRENT_DATE),
				updated_at = NOW()
			WHERE id = ?
		`, userID, name, profileType, description, id).Error
	}

	return tx.Exec(`
		INSERT INTO user_profiles (
			id, user_id, profile_type, is_primary, is_active, start_date, created_at, updated_at, name, description
		) VALUES (?, ?, ?, true, true, CURRENT_DATE, NOW(), NOW(), ?, ?)
	`, id, userID, profileType, name, description).Error
}

func upsertStudent(tx *gorm.DB, email, registration string) error {
	userID, err := getUserIDByEmail(tx, email)
	if err != nil {
		return err
	}

	var count int64
	if err := tx.Raw(`SELECT COUNT(1) FROM students WHERE user_id = ?`, userID).Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return tx.Exec(`
			UPDATE students
			SET status = 'active',
				registration_number = COALESCE(NULLIF(registration_number, ''), ?),
				updated_at = NOW()
			WHERE user_id = ?
		`, registration, userID).Error
	}

	return tx.Exec(`
		INSERT INTO students (user_id, registration_number, status, created_at, updated_at)
		VALUES (?, ?, 'active', NOW(), NOW())
	`, userID, registration).Error
}
