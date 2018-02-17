import {ADMIN} from '../../auth'

const ProjectForm = import ('./projects/Form')

export default {
  routes: [
    {
      path: "/donate/projects/edit/:id",
      component: ProjectForm
    }, {
      path: "/donate/projects/new",
      component: ProjectForm
    }, {
      path: "/donate/projects",
      component: import ('./projects/Index')
    }
  ],
  menus: [
    {
      icon: "pay-circle-o",
      label: "donate.dashboard.title",
      href: "donate",
      roles: [ADMIN],
      items: [
        {
          label: "donate.projects.index.title",
          href: "/donate/projects"
        }
      ]
    }
  ]
}
