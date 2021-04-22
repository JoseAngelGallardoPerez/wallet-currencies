package controllers

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-currencies/internal/abilities"
	"github.com/Confialink/wallet-currencies/internal/abilities/actions"
	"github.com/Confialink/wallet-currencies/internal/abilities/resources"
	"github.com/Confialink/wallet-currencies/internal/errcodes"
	"github.com/Confialink/wallet-currencies/internal/models"
	"github.com/Confialink/wallet-currencies/internal/serializers"
	"github.com/Confialink/wallet-currencies/internal/services/currencies"
	"github.com/Confialink/wallet-currencies/internal/services/files"
)

const mainCurrencyId = "main"

type Currencies struct {
	currencyService *currencies.Service
	currencyManager currencies.Manager
	filesService    *files.Service
	db              *gorm.DB
}

func NewCurrencies(currencyService *currencies.Service, currencyManager currencies.Manager, filesService *files.Service, db *gorm.DB) *Currencies {
	return &Currencies{currencyService: currencyService, currencyManager: currencyManager, filesService: filesService, db: db}
}

// AdminCreate creates a new currency
func (c *Currencies) AdminCreate(context *gin.Context) {
	if abilities.CanNot(context, actions.Create, resources.Currency) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	currency := currencies.Currency{}

	if err := context.ShouldBindJSON(&currency); err != nil {
		errors.AddShouldBindError(context, err)
		return
	}

	currencyModel, err := c.currencyManager.AddCurrency(currency)
	if err != nil {
		context.Error(err)
		return
	}

	serialized := serializers.AdminCurrency().Serialize(currencyModel)
	formCreatedResponse(context, serialized)
}

// Show returns currency by id or main currency
func (c *Currencies) Show(context *gin.Context) {
	if abilities.CanNot(context, actions.Read, resources.Currency) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	c.commonShow(context, serializers.Currency())
}

// Index returns list of all currencies
func (c *Currencies) Index(context *gin.Context) {
	if abilities.CanNot(context, actions.Index, resources.Currency) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	listParams := getListParams(context)
	if ok, validateErrors := listParams.Validate(); !ok {
		errcodes.AddErrorMeta(context, errcodes.BadCollectionParams, validateErrors)
		return
	}
	listParams.Pagination.PageSize = 0

	c.commonIndex(context, listParams, serializers.Currency())
}

// AdminShow returns currency by id or main currency for admin
func (c *Currencies) AdminShow(context *gin.Context) {
	if abilities.CanNot(context, actions.Read, resources.FullCurrency) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	c.commonShow(context, serializers.AdminCurrency())
}

// AdminIndex returns list of all currencies for admin
func (c *Currencies) AdminIndex(context *gin.Context) {
	if abilities.CanNot(context, actions.Index, resources.FullCurrency) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	listParams := getListParams(context)
	if ok, validateErrors := listParams.Validate(); !ok {
		errcodes.AddErrorMeta(context, errcodes.BadCollectionParams, validateErrors)
		return
	}

	c.commonIndex(context, listParams, serializers.AdminCurrency())
}

// AdminUpdateList handler to update list of currencies
func (c *Currencies) AdminUpdateList(context *gin.Context) {
	if abilities.CanNot(context, actions.Update, resources.FullCurrency) {
		errcodes.AddError(context, errcodes.Forbidden)
		return
	}

	var currenciesRequest = currencies.UpdateCurrenciesRequest{}

	if err := context.ShouldBindJSON(&currenciesRequest); err != nil {
		errors.AddShouldBindError(context, err)
		return
	}

	tx := c.db.Begin()
	currenciesService := c.currencyService.WrapContext(tx)

	if publicErrors := currenciesService.UpdateAdminCurrencies(&currenciesRequest); len(publicErrors) > 0 {
		tx.Rollback()
		errors.AddErrors(context, publicErrors...)
		context.Status(http.StatusBadRequest)
		return
	}

	tx.Commit()
	context.Status(http.StatusNoContent)
}

func (c *Currencies) GetLogo(context *gin.Context) {
	code := context.Params.ByName("code")
	currency, err := c.currencyService.FindByCode(code)
	if err != nil || currency.ID == 0 {
		privateErr := errors.PrivateError{
			Message:       "Currency not found",
			OriginalError: err,
		}
		privateErr.AddLogPair("currency code", code)
		errors.AddErrors(context, &privateErr)
		return
	}

	file, err := c.filesService.DownloadFileById(currency.LogoFileID)
	if err != nil {
		privateErr := errors.PrivateError{
			Message:       "Can't load a file",
			OriginalError: err,
		}
		errors.AddErrors(context, &privateErr)
		return
	}

	reader := bytes.NewReader(file.Data)
	context.DataFromReader(http.StatusOK, file.Size, file.ContentType, reader, map[string]string{})
}

func (c *Currencies) commonShow(context *gin.Context, serializer serializers.ICurrencySerializer) {
	currency, err := c.getCurrencyById(context.Param("id"))
	if err != nil {
		errors.AddErrors(context, &errors.PrivateError{OriginalError: err})
		return
	}
	if currency == nil {
		errcodes.AddError(context, errcodes.CurrencyNotFound)
	}

	formOkResponse(context, serializer.Serialize(currency))
}

func (c *Currencies) getCurrencyById(idParameter string) (*models.Currency, error) {
	if idParameter == mainCurrencyId {
		return c.currencyService.Main()
	} else {
		id, err := strconv.ParseUint(idParameter, 10, 32)
		if err != nil {
			return nil, err
		}
		convertedId := uint32(id)
		return c.currencyService.FindById(convertedId)
	}
}

func (c *Currencies) commonIndex(context *gin.Context, listParams *list_params.ListParams,
	serializer serializers.ICurrencySerializer) {

	currencyModels, err := c.currencyService.List(listParams)
	if err != nil {
		errors.AddErrors(context, &errors.PrivateError{OriginalError: err})
		return
	}

	r := serializer.SerializeList(currencyModels)

	formOkResponse(context, &r)
}
