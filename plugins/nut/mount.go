package nut

import (
	"fmt"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/spf13/viper"
)

func (p *Plugin) sitemap() ([]stm.URL, error) {
	var items []stm.URL
	for _, l := range viper.GetStringSlice("languages") {
		items = append(
			items,
			stm.URL{
				"loc": fmt.Sprintf("/?locale=%s", l),
			},
			stm.URL{
				"loc": fmt.Sprintf("/rss/%s", l),
			},
		)
	}
	return items, nil
}

// Mount register
func (p *Plugin) Mount() error {
	p.Sitemap.Register(p.sitemap)
	// --------------
	secret, err := p.secretKey()
	if err != nil {
		return err
	}

	rdr, err := NewHTMLRender(
		filepath.Join("themes", viper.GetString("server.theme")),
		p.renderFuncMap(),
	)
	if err != nil {
		return err
	}
	p.Router.HTMLRender = rdr
	im, err := p.I18n.Middleware()
	if err != nil {
		return err
	}
	store := sessions.NewCookieStore(secret)
	store.Options(sessions.Options{
		HttpOnly: true,
		Path:     "/",
		MaxAge:   60 * 20,
	})
	p.Router.Use(
		sessions.Sessions("_session_", store),
		im,
		p.Layout.CurrentUserMiddleware,
	)
	// --------------
	p.Router.GET("/", p.Layout.HTML("nut/index", p.getHome))

	api := p.Router.Group("/api")

	api.GET("/locales/:lang", p.Layout.JSON(p.getLocales))
	api.GET("/layout", p.Layout.JSON(p.getLayout))
	api.POST("/leave-words", p.Layout.JSON(p.createLeaveWord))

	ung := api.Group("/users")
	ung.POST("/sign-in", p.Layout.JSON(p.postUsersSignIn))
	ung.POST("/sign-up", p.Layout.JSON(p.postUsersSignUp))
	ung.POST("/confirm", p.Layout.JSON(p.postUsersConfirm))
	ung.POST("/unlock", p.Layout.JSON(p.postUsersUnlock))
	ung.POST("/forgot-password", p.Layout.JSON(p.postUsersForgotPassword))
	ung.POST("/reset-password", p.Layout.JSON(p.postUsersResetPassword))
	ung.GET("/confirm/:token", p.Layout.Redirect("/", p.getUsersConfirmToken))
	ung.GET("/unlock/:token", p.Layout.Redirect("/", p.getUsersUnlockToken))

	umg := api.Group("/users", p.Layout.MustSignInMiddleware)
	umg.GET("/logs", p.Layout.JSON(p.getUsersLogs))
	umg.GET("/profile", p.Layout.JSON(p.getUsersProfile))
	umg.POST("/profile", p.Layout.JSON(p.postUsersProfile))
	umg.POST("/change-password", p.Layout.JSON(p.postUsersChangePassword))
	umg.DELETE("/sign-out", p.Layout.JSON(p.deleteUsersSignOut))

	api.GET("/attachments", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexAttachments))
	api.POST("/attachments", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createAttachments))
	api.DELETE("/attachments/:id", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.destroyAttachments))

	ag := api.Group("/admin", p.Layout.MustAdminMiddleware)
	ag.GET("/site/status", p.Layout.JSON(p.getAdminSiteStatus))
	ag.POST("/site/info", p.Layout.JSON(p.postAdminSiteInfo))
	ag.POST("/site/author", p.Layout.JSON(p.postAdminSiteAuthor))
	ag.GET("/site/seo", p.Layout.JSON(p.getAdminSiteSeo))
	ag.POST("/site/seo", p.Layout.JSON(p.postAdminSiteSeo))
	ag.GET("/site/smtp", p.Layout.JSON(p.getAdminSiteSMTP))
	ag.POST("/site/smtp", p.Layout.JSON(p.postAdminSiteSMTP))
	ag.PATCH("/site/smtp", p.Layout.JSON(p.patchAdminSiteSMTP))
	ag.GET("/site/home", p.Layout.JSON(p.getAdminSiteHome))
	ag.POST("/site/home", p.Layout.JSON(p.postAdminSiteHome))
	ag.GET("/links", p.Layout.JSON(p.indexAdminLinks))
	ag.POST("/links", p.Layout.JSON(p.createAdminLink))
	ag.GET("/links/:id", p.Layout.JSON(p.showAdminLink))
	ag.POST("/links/:id", p.Layout.JSON(p.updateAdminLink))
	ag.DELETE("/links/:id", p.Layout.JSON(p.destroyAdminLink))
	ag.GET("/cards", p.Layout.JSON(p.indexAdminCards))
	ag.POST("/cards", p.Layout.JSON(p.createAdminCard))
	ag.GET("/cards/:id", p.Layout.JSON(p.showAdminCard))
	ag.POST("/cards/:id", p.Layout.JSON(p.updateAdminCard))
	ag.DELETE("/cards/:id", p.Layout.JSON(p.destroyAdminCard))
	ag.GET("/locales", p.Layout.JSON(p.indexAdminLocales))
	ag.POST("/locales", p.Layout.JSON(p.createAdminLocale))
	ag.GET("/locales/:id", p.Layout.JSON(p.showAdminLocale))
	ag.DELETE("/locales/:id", p.Layout.JSON(p.destroyAdminLocale))
	ag.GET("/friend-links", p.Layout.JSON(p.indexAdminFriendLinks))
	ag.POST("/friend-links", p.Layout.JSON(p.createAdminFriendLink))
	ag.GET("/friend-links/:id", p.Layout.JSON(p.showAdminFriendLink))
	ag.POST("/friend-links/:id", p.Layout.JSON(p.updateAdminFriendLink))
	ag.DELETE("/friend-links/:id", p.Layout.JSON(p.destroyAdminFriendLink))
	ag.GET("/leave-words", p.Layout.JSON(p.indexAdminLeaveWords))
	ag.DELETE("/leave-words/:id", p.Layout.JSON(p.destroyAdminLeaveWord))
	ag.GET("/users", p.Layout.JSON(p.indexAdminUsers))

	p.Queue.Register(SendEmailJob, p.doSendEmail)
	return nil
}
