package auth

import (
	"testing"
	"todo-service/src/infrastructure/authentication"
	"todo-service/src/infrastructure/delivery/graphql"
	"todo-service/src/infrastructure/storage"
	"todo-service/src/models"
	"todo-service/src/registry"
	"todo-service/tests/tools"
	"todo-service/utils"

	"github.com/google/uuid"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var router *echo.Echo
var db *gorm.DB

func AddUsersToDb() {
	users := []models.User{
		{
			ID:       uuid.MustParse("48f875c5-4d1f-4eb6-abbc-5e85dae826af"),
			Name:     "test_name",
			Email:    "test@gmail.com",
			Password: utils.HashPwd("12345"),
		},
		{
			ID:       uuid.MustParse("2fd3d635-3a46-4085-a854-81f76a25cfd0"),
			Name:     "test_name2",
			Email:    "test2@gmail.com",
			Password: utils.HashPwd("12345"),
		},
	}
	result := db.Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}
}

var _ = BeforeSuite(func() {
	viper.AddConfigPath("../../../conf")
	viper.SetConfigName("test_config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	db = storage.InitPostgres(logger)

	jc := authentication.NewJwtConfigurator(logger, "../../../rsa_keys/private_key.pem",
		"../../../rsa_keys/public_key.pem")

	// Register and create controller
	useCase := registry.NewRegistry(db, jc).NewUseCase()

	router = echo.New()

	// Initialize Echo instance
	graphql.NewGraphqlRouter(router, useCase)
})

var _ = Describe("Auth", func() {
	type signUp struct {
		Auth struct {
			Data struct {
				IsCreated bool `json:"isCreated"`
			} `graphql:"signUp(input: {name:$name, email:$email, password:$password})"`
		} `json:"auth"`
	}

	type signIn struct {
		Auth struct {
			Data struct {
				AccessToken  string `json:"accessToken"`
				RefreshToken string `json:"refreshToken"`
			} `graphql:"signIn(email:$email, password:$password)"`
		} `json:"auth"`
	}

	Describe("Sign up", func() {
		BeforeEach(func() {
			err := db.Migrator().DropTable(&models.Todo{})
			if err != nil {
				panic(err)
			}

			err = db.Migrator().DropTable(&models.User{})
			if err != nil {
				panic(err)
			}

			err = db.AutoMigrate(&models.User{})
			if err != nil {
				panic(err)
			}
			err = db.AutoMigrate(&models.Todo{})
			if err != nil {
				panic(err)
			}
		})

		Context("With valid parameters", func() {
			It("success create user", func() {
				var q signUp
				var wantQ signUp

				wantQ.Auth.Data.IsCreated = true

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    "test@email.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))
			})
		})

		Context("With already existing email", func() {
			It("returns err about invalid fields", func() {
				var q signUp
				var wantQ signUp
				wantQ.Auth.Data.IsCreated = true

				var q2 signUp
				var wantQ2 signUp

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    "test@email.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

				err = tools.DoMutate(&q2, variables, "", router)
				Expect(err.Error()).To(Equal("Message: email already exist, Locations: [], Extensions: map[]"))
				Expect(q2).To(Equal(wantQ2))
			})
		})

		Context("With different name", func() {
			It("returns err about invalid field", func() {
				var q signUp
				var wantQ signUp

				variables := map[string]interface{}{
					"name":     "t",
					"email":    "test@email.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: Parameters incorrectly formatted or out of range (Name), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})

			It("returns err about invalid field", func() {
				var q signUp
				var wantQ signUp

				variables := map[string]interface{}{
					"name":     tools.GenerateRandomString(129),
					"email":    "test@email.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: Parameters incorrectly formatted or out of range (Name), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})

			It("returns success message", func() {
				var q signUp
				var wantQ signUp
				wantQ.Auth.Data.IsCreated = true

				variables := map[string]interface{}{
					"name":     tools.GenerateRandomString(128),
					"email":    "test@email.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))
			})
		})

		Context("With different email", func() {
			It("returns err about invalid field", func() {
				var q signUp
				var wantQ signUp

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    "testemail.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: Parameters incorrectly formatted or out of range (Email), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})

			It("returns err about invalid field", func() {
				var q signUp
				var wantQ signUp

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    tools.GenerateRandomString(246) + "@email.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: Parameters incorrectly formatted or out of range (Email), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})

			It("returns success message", func() {
				var q signUp
				var wantQ signUp
				wantQ.Auth.Data.IsCreated = true

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    tools.GenerateRandomString(245) + "@email.com",
					"password": "1234567",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))
			})
		})

		Context("With different password", func() {
			It("returns err about invalid field", func() {
				var q signUp
				var wantQ signUp

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    "teste@mail.com",
					"password": "12345",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: Parameters incorrectly formatted or out of range (Password), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})

			It("returns err about invalid field", func() {
				var q signUp
				var wantQ signUp

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    "teste@mail.com",
					"password": tools.GenerateRandomString(65),
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: Parameters incorrectly formatted or out of range (Password), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})

			It("returns success message", func() {
				var q signUp
				var wantQ signUp
				wantQ.Auth.Data.IsCreated = true

				variables := map[string]interface{}{
					"name":     "test_name",
					"email":    "teste@mail.com",
					"password": tools.GenerateRandomString(64),
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))
			})
		})
	})

	Describe("Sign in", func() {
		BeforeEach(func() {
			err := db.Migrator().DropTable(&models.Todo{})
			if err != nil {
				panic(err)
			}

			err = db.Migrator().DropTable(&models.User{})
			if err != nil {
				panic(err)
			}

			err = db.AutoMigrate(&models.User{})
			if err != nil {
				panic(err)
			}
			err = db.AutoMigrate(&models.Todo{})
			if err != nil {
				panic(err)
			}

			AddUsersToDb()

		})

		Context("With valid parameters", func() {
			It("success create user", func() {
				var q signIn

				variables := map[string]interface{}{
					"email":    "test@gmail.com",
					"password": "12345",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).To(BeNil())
				Expect(q.Auth.Data.AccessToken).ToNot(Equal(""))
				Expect(q.Auth.Data.RefreshToken).ToNot(Equal(""))
			})
		})

		Context("Doesn't exist email", func() {
			It("returns err about email not found", func() {
				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    "test4@gmail.com",
					"password": "123455",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: email not found, Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})
		})

		Context("With invalid password", func() {
			It("returns err about invalid password", func() {
				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    "test@gmail.com",
					"password": "123455",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).To(Equal("Message: invalid password, Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))
			})
		})

		Context("With invalid name of param (email)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email1":   "test@gmail.com",
					"password": "123455",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid type of email (int)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    1,
					"password": "123455",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid type of email (bool)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    true,
					"password": "123455",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid type of email  (float)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    1.5,
					"password": "123455",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid name of param (password)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":     "test@gmail.com",
					"password1": "123455",
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid type of password (int)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    "test@gmail.com",
					"password": 1,
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid type of password (bool)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    "test@gmail.com",
					"password": true,
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid type of password  (float)", func() {
			It("error: error type", func() {

				var q signIn
				var wantQ signIn

				variables := map[string]interface{}{
					"email":    "test@gmail.com",
					"password": 1.5,
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err.Error()).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})
	})
})
