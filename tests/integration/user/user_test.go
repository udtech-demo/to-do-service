package user

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

var signInUser1Resp userInfo

type userInfo struct {
	User  models.User
	Todos []models.Todo
	signIn
}

type signIn struct {
	Auth struct {
		Data struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		} `graphql:"signIn(email:$email, password:$password)"`
	} `json:"auth"`
}

type me struct {
	Data struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	} `graphql:"me"`
}

func initMeStruct(info userInfo) me {
	return me{
		Data: struct {
			ID    uuid.UUID `json:"id"`
			Name  string    `json:"name"`
			Email string    `json:"email"`
		}{
			ID:    info.User.ID,
			Name:  info.User.Name,
			Email: info.User.Email,
		},
	}
}

func TestTodo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Suite")
}

var router *echo.Echo
var db *gorm.DB

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

func SignIn(email, pwd string) (signIn, error) {
	var q signIn
	variables := map[string]interface{}{
		"email":    email,
		"password": pwd,
	}

	err := tools.DoMutate(&q, variables, "", router)
	return q, err
}

func AddUsersToDb() {
	users := []models.User{
		{
			ID:       uuid.MustParse("48f875c5-4d1f-4eb6-abbc-5e85dae826af"),
			Name:     "test_name",
			Email:    "test@gmail.com",
			Password: utils.HashPwd("12345"),
		},
	}
	result := db.Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}

	signInUser1Resp.User = users[0]
	signInUser1Resp.User.Password = "12345"

}

var _ = Describe("User", func() {

	Describe("Me", func() {
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

			signInUser1Resp.signIn, err = SignIn(signInUser1Resp.User.Email, signInUser1Resp.User.Password)
			if err != nil {
				panic(err)
			}

		})

		Context("With valid access token", func() {
			It("success get user information", func() {

				var q me

				wantQ := initMeStruct(signInUser1Resp)

				err := tools.DoQuery(&q, nil, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))
			})
		})

		Context("With invalid access token", func() {
			It("error: 401 Unauthorized", func() {

				var q me
				var wantQ me

				err := tools.DoQuery(&q, nil, signInUser1Resp.Auth.Data.AccessToken+"r", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("Without access token", func() {
			It("error: Access Denied", func() {

				var q me
				var wantQ me

				err := tools.DoQuery(&q, nil, "", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})
	})
})
