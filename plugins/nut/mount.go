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

	ug := p.Router.Group("/users")
	ug.POST("/sign-in", p.postUsersSignIn)
	ug.POST("/sign-up", p.postUsersSignUp)
	ug.POST("/confirm", p.postUsersConfirm)
	ug.POST("/unlock", p.postUsersUnlock)
	ug.POST("/forgot-password", p.postUsersForgotPassword)
	ug.POST("/reset-password", p.postUsersResetPassword)
	ug.GET("/confirm/{token}", p.getUsersConfirmToken)
	ug.GET("/unlock/{token}", p.getUsersUnlockToken)
	ug.GET("/logs", p.Layout.MustSignInMiddleware, p.getUsersLogs)
	ug.GET("/profile", p.Layout.MustSignInMiddleware, p.getUsersProfile)
	ug.POST("/profile", p.Layout.MustSignInMiddleware, p.postUsersProfile)
	ug.POST("/change-password", p.Layout.MustSignInMiddleware, p.postUsersChangePassword)
	ug.DELETE("/sign-out", p.Layout.MustSignInMiddleware, p.deleteUsersSignOut)

	atg := p.Router.Group("/attachments", p.Layout.MustSignInMiddleware)
	atg.GET("/", p.indexAttachments)
	atg.POST("/attachments", p.createAttachments)
	atg.DELETE("/attachments/{id}", p.destroyAttachments)

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
