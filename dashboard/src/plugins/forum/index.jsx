import {USER, ADMIN} from '../../auth'

const FormTag = import ('./tags/Form')
const FormCatalog = import ('./catalogs/Form')
const FormTopic = import ('./topics/Form')
const FormPost = import ('./posts/Form')

export default {
  menus: [
    {
      icon: "tablet",
      label: "forum.dashboard.title",
      href: "forum",
      roles: [
        USER, ADMIN
      ],
      items: [
        {
          label: "forum.topics.index.title",
          href: "/forum/topics"
        }, {
          label: "forum.posts.index.title",
          href: "/forum/posts"
        }, {
          label: "forum.tags.index.title",
          href: "/forum/tags",
          roles: [ADMIN]
        }, {
          label: "forum.catalogs.index.title",
          href: "/forum/catalogs",
          roles: [ADMIN]
        }
      ]
    }
  ],
  routes: [
    {
      path: "/forum/catalogs/edit/:id",
      component: FormCatalog
    }, {
      path: "/forum/catalogs/new",
      component: FormCatalog
    }, {
      path: "/forum/catalogs",
      component: import ('./catalogs/Index')
    }, {
      path: "/forum/tags/edit/:id",
      component: FormTag
    }, {
      path: "/forum/tags/new",
      component: FormTag
    }, {
      path: "/forum/tags",
      component: import ('./tags/Index')
    }, {
      path: "/forum/topics/edit/:id",
      component: FormTopic
    }, {
      path: "/forum/topics/new",
      component: FormTopic
    }, {
      path: "/forum/topics",
      component: import ('./topics/Index')
    }, {
      path: "/forum/posts/edit/:id",
      component: FormPost
    }, {
      path: "/forum/posts/new/:topicId",
      component: FormPost
    }, {
      path: "/forum/posts",
      component: import ('./posts/Index')
    }
  ]
}
