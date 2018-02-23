#include "server.h"

namespace arche {

void HTTPTimeServer::handleHelp(const std::string &name,
                                const std::string &value) {
  HelpFormatter helpFormatter(options());
  helpFormatter.setCommand(commandName());
  helpFormatter.setUsage("OPTIONS");
  helpFormatter.setHeader(
      "A web server that serves the current date and time.");
  helpFormatter.format(std::cout);
  stopOptionsProcessing();
  _helpRequested = true;
}

int HTTPTimeServer::main(const std::vector<std::string> &args) {
  if (!_helpRequested) {
    unsigned short port =
        (unsigned short)config().getInt("HTTPTimeServer.port", 8080);
    std::string format(config().getString("HTTPTimeServer.format",
                                          DateTimeFormat::SORTABLE_FORMAT));

    ServerSocket svs(port);
    HTTPServer srv(new TimeRequestHandlerFactory(format), svs,
                   new HTTPServerParams);
    srv.start();
    waitForTerminationRequest();
    srv.stop();
  }
  return Application::EXIT_OK;
}

void HTTPTimeServer::defineOptions(OptionSet &options) {
  ServerApplication::defineOptions(options);

  options.addOption(Option("help", "h", "display argument help information")
                        .required(false)
                        .repeatable(false)
                        .callback(OptionCallback<HTTPTimeServer>(
                            this, &HTTPTimeServer::handleHelp)));
}

} // namespace arche
