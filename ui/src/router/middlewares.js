/**
 * Current user state initialization
 * @WARN Must be always first in middleware chain
 */
export function initCurrentUserStateMiddleware (to, from, next) {
  next()
}

/**
 * Check access permission to auth routes
 */
export function checkAccessMiddleware (to, from, next) {
  next()
}

export function setPageTitleMiddleware (to, from, next) {
  const pageTitle = to.matched.find(item => item.meta.title)

  if (pageTitle) window.document.title = pageTitle.meta.title
  next()
}
