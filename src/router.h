#ifndef ARCHE_ROUTER_H_
#define ARCHE_ROUTER_H_

#include "common.h"
#include "forum.h"

using Poco::Net::HTTPRequestHandler;
using Poco::Net::HTTPRequestHandlerFactory;
using Poco::Net::HTTPServerRequest;

namespace arche {

class TimeRequestHandlerFactory : public HTTPRequestHandlerFactory {
public:
  TimeRequestHandlerFactory(const std::string &format) : _format(format) {}
  HTTPRequestHandler *createRequestHandler(const HTTPServerRequest &request);

private:
  std::string _format;
};

} // namespace arche

#endif
