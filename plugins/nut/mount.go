package nut

import (
	"github.com/chonglou/arche/web/queue"
)

// Mount register
func (p *Plugin) Mount() error {
	p.Router.Use(p.Layout.CurrentUserMiddleware)

	p.Router.GET("/layout", p.getLayout)
	p.Router.GET("/locales/{lang}", p.getLocales)
	p.Router.POST("/install", p.postInstall)
	p.Router.POST("/leave-words", p.createLeaveWord)

	ung := p.Router.Group("/users")
	ung.POST("/sign-in", p.postUsersSignIn)
	ung.POST("/sign-up", p.postUsersSignUp)
	ung.POST("/confirm", p.postUsersConfirm)
	ung.POST("/unlock", p.postUsersUnlock)
	ung.POST("/forgot-password", p.postUsersForgotPassword)
	ung.POST("/reset-password", p.postUsersResetPassword)
	ung.GET("/confirm/{token}", p.getUsersConfirmToken)
	ung.GET("/unlock/{token}", p.getUsersUnlockToken)
	umg := p.Router.Group("/users", p.Layout.MustSignInMiddleware)
	umg.GET("/logs", p.getUsersLogs)
	umg.GET("/profile", p.getUsersProfile)
	umg.POST("/profile", p.postUsersProfile)
	umg.POST("/change-password", p.postUsersChangePassword)
	umg.DELETE("/sign-out", p.deleteUsersSignOut)

	atg := p.Router.Group("/attachments", p.Layout.MustSignInMiddleware)
	atg.GET("/", p.indexAttachments)
	atg.POST("/", p.createAttachments)
	atg.DELETE("/{id}", p.destroyAttachments)

	ag := p.Router.Group("/admin", p.Layout.MustAdminMiddleware)
	ag.GET("/site/status", p.getAdminSiteStatus)
	ag.DELETE("/site/clear-cache", p.deleteAdminSiteClearCache)
	ag.POST("/site/info", p.postAdminSiteInfo)
	ag.POST("/site/author", p.postAdminSiteAuthor)
	ag.GET("/site/seo", p.getAdminSiteSeo)
	ag.POST("/site/seo", p.postAdminSiteSeo)
	ag.GET("/site/smtp", p.getAdminSiteSMTP)
	ag.POST("/site/smtp", p.postAdminSiteSMTP)
	ag.PATCH("/site/smtp", p.patchAdminSiteSMTP)
	ag.GET("/site/home", p.getAdminSiteHome)
	ag.POST("/site/home", p.postAdminSiteHome)
	ag.GET("/links", p.indexAdminLinks)
	ag.POST("/links", p.createAdminLink)
	ag.GET("/links/{id}", p.showAdminLink)
	ag.POST("/links/{id}", p.updateAdminLink)
	ag.DELETE("/links/{id}", p.destroyAdminLink)
	ag.GET("/cards", p.indexAdminCards)
	ag.POST("/cards", p.createAdminCard)
	ag.GET("/cards/{id}", p.showAdminCard)
	ag.POST("/cards/{id}", p.updateAdminCard)
	ag.DELETE("/cards/{id}", p.destroyAdminCard)
	ag.GET("/locales", p.indexAdminLocales)
	ag.POST("/locales", p.createAdminLocale)
	ag.GET("/locales/{id}", p.showAdminLocale)
	ag.DELETE("/locales/{id}", p.destroyAdminLocale)
	ag.GET("/friend-links", p.indexAdminFriendLinks)
	ag.POST("/friend-links", p.createAdminFriendLink)
	ag.GET("/friend-links/{id}", p.showAdminFriendLink)
	ag.POST("/friend-links/{id}", p.updateAdminFriendLink)
	ag.DELETE("/friend-links/{id}", p.destroyAdminFriendLink)
	ag.GET("/leave-words", p.indexAdminLeaveWords)
	ag.DELETE("/leave-words/{id}", p.destroyAdminLeaveWord)
	ag.GET("/users", p.indexAdminUsers)

	queue.Register(SendEmailJob, p.doSendEmail)
	return nil
}
