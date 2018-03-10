package forum

// Mount register
func (p *Plugin) Mount() error {
	// ---------------
	api := p.Router.Group("/api/forum")

	api.GET("/catalogs", p.Layout.JSON(p.indexCatalogs))
	api.POST("/catalogs", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.createCatalog))
	api.GET("/catalogs/:id", p.Layout.JSON(p.showCatalog))
	api.POST("/catalogs/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.updateCatalog))
	api.DELETE("/catalogs/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyCatalog))

	api.GET("/tags", p.Layout.JSON(p.indexTags))
	api.POST("/tags", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.createTag))
	api.GET("/tags/:id", p.Layout.JSON(p.showTag))
	api.POST("/tags/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.updateTag))
	api.DELETE("/tags/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyTag))

	api.GET("/topics", p.Layout.JSON(p.indexTopics))
	api.POST("/topics", p.canEditTopic, p.Layout.JSON(p.createTopic))
	api.GET("/topics/:id", p.Layout.JSON(p.showTopic))
	api.POST("/topics/:id", p.canEditTopic, p.Layout.JSON(p.updateTopic))
	api.DELETE("/topics/:id", p.canEditTopic, p.Layout.JSON(p.destroyTopic))

	api.GET("/posts", p.Layout.JSON(p.indexPosts))
	api.POST("/posts", p.canEditPost, p.Layout.JSON(p.createPost))
	api.GET("/posts/:id", p.Layout.JSON(p.showPost))
	api.POST("/posts/:id", p.canEditPost, p.Layout.JSON(p.updatePost))
	api.DELETE("/posts/:id", p.canEditPost, p.Layout.JSON(p.destroyPost))
	return nil
}
