#include "router.h"
namespace arche {

HTTPRequestHandler *
Router::createRequestHandler(const HTTPServerRequest &request) {
  if (request.getURI() == "/")
    return new forum::TimeRequestHandler(_format);
  else
    return 0;
}
} // namespace arche
