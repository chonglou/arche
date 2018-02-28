import React from 'react'

import {USER, ADMIN} from '../../auth'

import Home from './Home'
import AttachmentsIndex from './attachments/Index'
import AdminCardsIndex from './admin/cards/Index'
import AdminCardsForm from './admin/cards/Form'
import AdminLinksIndex from './admin/links/Index'
import AdminLinksForm from './admin/links/Form'
import AdminFriendLinksIndex from './admin/friend-links/Index'
import AdminFriendLinksForm from './admin/friend-links/Form'
import AdminLocalesIndex from './admin/locales/Index'
import AdminLocalesForm from './admin/locales/Form'
import AdminLeaveWordsIndex from './admin/leave-words/Index'
import AdminUsersIndex from './admin/users/Index'
import AdminSiteHome from './admin/site/Home'
import AdminSiteSmtp from './admin/site/Smtp'
import AdminSiteSeo from './admin/site/Seo'
import AdminSiteAuthor from './admin/site/Author'
import AdminSiteInfo from './admin/site/Info'
import AdminSiteStatus from './admin/site/Status'

import LeaveWordsNew from './leave-words/New'
import UsersChangePassword from './users/ChangePassword'
import UsersProfile from './users/Profile'
import UsersLogs from './users/Logs'
import UsersResetPassword from './users/ResetPassword'
import UsersSignIn from './users/SignIn'
import UsersSignUp from './users/SignUp'
import UsersEmailForm from './users/EmailForm'

const UsersConfirm = () => (<UsersEmailForm action="confirm"/>)
const UsersUnlock = () => (<UsersEmailForm action="unlock"/>)
const UsersForgotPassword = () => (<UsersEmailForm action="forgot-password"/>)

export default {
  routes: [
    {
      path: "/",
      component: Home
    }, {
      path: "/users/sign-in",
      component: UsersSignIn
    }, {
      path: "/users/sign-up",
      component: UsersSignUp
    }, {
      path: "/users/confirm",
      component: UsersConfirm
    }, {
      path: "/users/unlock",
      component: UsersUnlock
    }, {
      path: "/users/forgot-password",
      component: UsersForgotPassword
    }, {
      path: "/users/reset-password/:token",
      component: UsersResetPassword
    }, {
      path: "/users/logs",
      component: UsersLogs
    }, {
      path: "/users/profile",
      component: UsersProfile
    }, {
      path: "/users/change-password",
      component: UsersChangePassword
    }, {
      path: "/leave-words/new",
      component: LeaveWordsNew
    }, {
      path: "/admin/site/status",
      component: AdminSiteStatus
    }, {
      path: "/admin/site/info",
      component: AdminSiteInfo
    }, {
      path: "/admin/site/author",
      component: AdminSiteAuthor
    }, {
      path: "/admin/site/seo",
      component: AdminSiteSeo
    }, {
      path: "/admin/site/smtp",
      component: AdminSiteSmtp
    }, {
      path: "/admin/site/home",
      component: AdminSiteHome
    }, {
      path: "/admin/users",
      component: AdminUsersIndex
    }, {
      path: "/admin/leave-words",
      component: AdminLeaveWordsIndex
    }, {
      path: "/admin/locales/edit/:id",
      component: AdminLocalesForm
    }, {
      path: "/admin/locales/new",
      component: AdminLocalesForm
    }, {
      path: "/admin/locales",
      component: AdminLocalesIndex
    }, {
      path: "/admin/friend-links/edit/:id",
      component: AdminFriendLinksForm
    }, {
      path: "/admin/friend-links/new",
      component: AdminFriendLinksForm
    }, {
      path: "/admin/friend-links",
      component: AdminFriendLinksIndex
    }, {
      path: "/admin/links/edit/:id",
      component: AdminLinksForm
    }, {
      path: "/admin/links/new",
      component: AdminLinksForm
    }, {
      path: "/admin/links",
      component: AdminLinksIndex
    }, {
      path: "/admin/cards/edit/:id",
      component: AdminCardsForm
    }, {
      path: "/admin/cards/new",
      component: AdminCardsForm
    }, {
      path: "/admin/cards",
      component: AdminCardsIndex
    }, {
      path: "/attachments",
      component: AttachmentsIndex
    }
  ],
  menus: [
    {
      icon: "user",
      label: "nut.self.title",
      href: "personal",
      roles: [
        USER, ADMIN
      ],
      items: [
        {
          label: "nut.users.logs.title",
          href: "/users/logs"
        }, {
          label: "nut.users.profile.title",
          href: "/users/profile"
        }, {
          label: "nut.users.change-password.title",
          href: "/users/change-password"
        }, {
          label: "nut.attachments.index.title",
          href: "/attachments"
        }
      ]
    }, {
      icon: "setting",
      label: "nut.settings.title",
      href: "settings",
      roles: [ADMIN],
      items: [
        {
          label: "nut.admin.site.status.title",
          href: "/admin/site/status"
        }, {
          label: "nut.admin.site.info.title",
          href: "/admin/site/info"
        }, {
          label: "nut.admin.site.author.title",
          href: "/admin/site/author"
        }, {
          label: "nut.admin.site.seo.title",
          href: "/admin/site/seo"
        }, {
          label: "nut.admin.site.smtp.title",
          href: "/admin/site/smtp"
        }, {
          label: "nut.admin.site.donate.title",
          href: "/admin/site/donate"
        }, {
          label: "nut.admin.site.home.title",
          href: "/admin/site/home"
        }, {
          label: "nut.admin.links.index.title",
          href: "/admin/links"
        }, {
          label: "nut.admin.cards.index.title",
          href: "/admin/cards"
        }, {
          label: "nut.admin.locales.index.title",
          href: "/admin/locales"
        }, {
          label: "nut.admin.friend-links.index.title",
          href: "/admin/friend-links"
        }, {
          label: "nut.admin.leave-words.index.title",
          href: "/admin/leave-words"
        }, {
          label: "nut.admin.users.index.title",
          href: "/admin/users"
        }
      ]
    }
  ]
}
