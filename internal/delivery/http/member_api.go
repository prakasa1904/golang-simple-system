package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MemberAPIController struct {
	log       *logrus.Logger
	myUsecase *member.UseCase
	view      *render.Engine
}

func NewMemberAPIController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
) *MemberAPIController {
	// init module repositories
	myRepository := member.NewRepository(log)

	// init module usecase
	myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &MemberAPIController{
		log:       log,
		myUsecase: myUsecase,
		view:      view,
	}
}

func (c *MemberAPIController) Add(w http.ResponseWriter, r *http.Request) {
	var payloadRegister = new(member.RegisterRequest)
	err := json.NewDecoder(r.Body).Decode(&payloadRegister)
	if err != nil {
		c.log.Warnf("Json decoder error : %+v", err)
	}

	log.Println("payloadRegisterpayloadRegisterpayloadRegister")
	log.Println(payloadRegister)
	log.Println("payloadRegisterpayloadRegisterpayloadRegister")

	currentURL, err := url.Parse(payloadRegister.Request.CurrentURL)
	if err != nil {
		c.log.Warnf("Parse payload current url error : %+v", err)
	}

	newUser, err := c.myUsecase.Create(r.Context(), payloadRegister)
	if err != nil {
		c.log.Warnf("Find users error : %+v", err)
	}

	if newUser.ID == 0 {
		c.log.Warnf("Failed to create new user")
	}

	filter := member.ConvertQueryToFilter(currentURL)
	limit := member.ConvertQueryToLimit(currentURL)
	order := member.ConvertQueryToOrder(currentURL)

	users, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find users error : %+v", err)
	}

	c.view.Set("users", users)

	err = c.view.HTML(w).RenderClean("views/pages/home/list-members.html")
	if err != nil {
		c.log.Warnf("RenderClean error : %+v", err)
	}

}
