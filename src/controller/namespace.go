package controller

import (
	"github.com/labstack/echo"
	"github.com/tidwall/gjson"
	"github.com/two-rabbits/ranvier/src/model"
	"github.com/two-rabbits/ranvier/src/service"
	"io/ioutil"
	"net/http"
	"strconv"
)

const NamespaceControllerKey = "NamespaceController"

type NamespaceController interface {
	StratifiedLevelConfigController
}

type namespaceControllerImpl struct {
	LevelConfigService service.NamespaceLevelConfigService `inject:"NamespaceLevelConfigService"`
}

func (n *namespaceControllerImpl) GetAll(c echo.Context) error {
	cluster := c.Param("cluster")

	if cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	resp, err := n.LevelConfigService.GetAll(cluster)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (n *namespaceControllerImpl) Get(c echo.Context) error {
	panic("implement me")
}

func (n *namespaceControllerImpl) MergedQuery(c echo.Context) error {
	key := c.QueryParam("key")
	namespace := c.Param("namespace")
	cluster := c.Param("cluster")

	if namespace == "" || cluster == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	data, err := n.LevelConfigService.MergedQuery(cluster, namespace, key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (n *namespaceControllerImpl) Update(c echo.Context) error {
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")

	if namespace == "" || cluster == "" {
		return c.NoContent(http.StatusNotFound)
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	appConfig, err := n.LevelConfigService.Update(cluster, namespace, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, appConfig)
}

func (n *namespaceControllerImpl) Rollback(c echo.Context) error {
	version := c.Param("version")
	namespace := c.Param("namespace")
	cluster := c.Param("cluster")

	if namespace == "" || cluster == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	ver, err := strconv.Atoi(version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "version must be an int")
	}

	config, err := n.LevelConfigService.Rollback(cluster, namespace, ver)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, config)
}

func (n *namespaceControllerImpl) Query(c echo.Context) error {
	key := c.QueryParam("key")
	namespace := c.Param("namespace")
	cluster := c.Param("cluster")

	if namespace == "" || cluster == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	data, err := n.LevelConfigService.Query(cluster, namespace, key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func (n *namespaceControllerImpl) Create(c echo.Context) error {
	namespace := c.Param("namespace")
	cluster := c.Param("cluster")

	if namespace == "" || cluster == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !gjson.ValidBytes(data) {
		return c.NoContent(http.StatusBadRequest)
	}

	namespaceConfig, err := n.LevelConfigService.Create(cluster, namespace, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, namespaceConfig)
}

func (n *namespaceControllerImpl) GetRoutes() []model.Route {
	return []model.Route{
		model.NewRoute(http.MethodPost, "/config/:cluster/:namespace", true, n.Create),
		model.NewRoute(http.MethodPut, "/config/:cluster/:namespace/rollback/:version", true, n.Rollback),
		model.NewRoute(http.MethodPut, "/config/:cluster/:namespace", true, n.Update),
		model.NewRoute(http.MethodGet, "/config/:cluster/:namespace", true, n.Query),
		model.NewRoute(http.MethodGet, "/config/:cluster/namespaces", true, n.GetAll),
		model.NewRoute(http.MethodGet, "/config/:cluster/namespaces/:name", true, n.Get),
		model.NewRoute(http.MethodGet, "/config/:cluster/:namespace", false, n.MergedQuery),
	}
}
