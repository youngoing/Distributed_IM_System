package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

// registerServiceHandler 处理注册服务的请求
func RegisterServiceHandler(c *gin.Context) {
	var service ConsulService
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 使用 Consul API 注册服务
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Consul client"})
		return
	}

	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Check: &api.AgentServiceCheck{
			HTTP:     service.Check.HTTP,
			Interval: service.Check.Interval,
			Timeout:  service.Check.Timeout,
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service registered successfully"})
}

// deregisterServiceHandler 处理注销服务的请求
func DeregisterServiceHandler(c *gin.Context) {
	serviceID := c.Param("id")

	// 使用 Consul API 注销服务
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Consul client"})
		return
	}

	err = client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deregister service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deregistered successfully"})
}

// listServicesHandler 处理查询所有服务的请求
func ListServicesHandler(c *gin.Context) {
	// 使用 Consul API 获取所有服务信息
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Consul client"})
		return
	}

	services, err := client.Agent().Services()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

// updateServiceHandler 处理更新服务的请求
func UpdateServiceHandler(c *gin.Context) {
	serviceID := c.Param("id")
	var service ConsulService

	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 先注销旧服务
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Consul client"})
		return
	}

	err = client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deregister old service"})
		return
	}

	// 再注册新服务
	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Check: &api.AgentServiceCheck{
			HTTP:     service.Check.HTTP,
			Interval: service.Check.Interval,
			Timeout:  service.Check.Timeout,
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register new service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
}
