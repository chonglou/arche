#ifndef ARCHE_SERVER_H_
#define ARCHE_SERVER_H_

#include "common.h"
#include "router.h"

using Poco::DateTimeFormat;
using Poco::Net::HTTPServer;
using Poco::Net::HTTPServerParams;
using Poco::Net::ServerSocket;
using Poco::ThreadPool;
using Poco::Util::HelpFormatter;
using Poco::Util::Option;
using Poco::Util::OptionCallback;
using Poco::Util::OptionSet;
using Poco::Util::ServerApplication;

namespace arche {

class HTTPTimeServer : public Poco::Util::ServerApplication {
public:
  HTTPTimeServer() : _helpRequested(false) {}

  ~HTTPTimeServer() {}

protected:
  void initialize(Application &self) {
    loadConfiguration();
    ServerApplication::initialize(self);
  }

  void uninitialize() { ServerApplication::uninitialize(); }

  void handleHelp(const std::string &name, const std::string &value);
  int main(const std::vector<std::string> &args);
  void defineOptions(OptionSet &options);

private:
  bool _helpRequested;
};

} // namespace arche

#endif
