package handlers

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/Vladroon22/TestTask-ITK-Academy/internal/entity"
	"github.com/Vladroon22/TestTask-ITK-Academy/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	ErrWrongInput = errors.New("incorrect input of data")
)

type Handlers struct {
	cache map[string]entity.WalletData
	srv   service.Servicer
	mu    sync.RWMutex
}

func NewHandler(s service.Servicer) *Handlers {
	return &Handlers{srv: s, cache: make(map[string]entity.WalletData), mu: sync.RWMutex{}}
}

func (h *Handlers) WalletOperation(c *gin.Context) {
	wallet := entity.WalletData{}

	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrWrongInput})
		log.Println(err)
		return
	}

	validW, errV := entity.Validate(wallet)
	if errV != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errV.Error()})
		log.Println(errV)
		return
	}

	if err := h.srv.WalletOperation(c.Request.Context(), validW); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, "Your finacial operation is successfull")
}

func (h *Handlers) GetBalance(c *gin.Context) {
	UUID := c.Param("id")

	h.mu.RLock()
	cacheWallet, ok := h.cache[UUID]
	h.mu.RUnlock()
	if ok {
		c.JSON(http.StatusOK, cacheWallet)
		return
	}

	wallet, err := h.srv.GetBalance(c.Request.Context(), UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, wallet)

	h.mu.Lock()
	h.cache[UUID] = wallet
	h.mu.Unlock()
}
