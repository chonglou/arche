#ifndef ARCHE_FORUM_H_
#define ARCHE_FORUM_H_

#include "common.h"

using Poco::DateTimeFormatter;
using Poco::Net::HTTPRequestHandler;
using Poco::Net::HTTPServerRequest;
using Poco::Net::HTTPServerResponse;
using Poco::Timestamp;
using Poco::Util::Application;

namespace arche {
namespace forum {

class TimeRequestHandler : public HTTPRequestHandler {
public:
  TimeRequestHandler(const std::string &format) : _format(format) {}
  void handleRequest(HTTPServerRequest &request, HTTPServerResponse &response);

private:
  std::string _format;
};

} // namespace forum
} // namespace arche

#endif
