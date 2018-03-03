import ProjectsForm from './projects/Form'
import ProjectsIndex from './projects/Index'

import {USER, ADMIN} from '../../auth'

export default {
  routes: [
    {
      path: "/donate/projects/edit/:id",
      component: ProjectsForm
    }, {
      path: "/donate/projects/new",
      component: ProjectsForm
    }, {
      path: "/donate/projects",
      component: ProjectsIndex
    }
  ],
  menus: [
    {
      icon: "pay-circle-o",
      label: "donate.dashboard.title",
      href: "donate",
      roles: [
        USER, ADMIN
      ],
      items: [
        {
          label: "donate.projects.index.title",
          href: "/donate/projects"
        }
      ]
    }
  ]
}
