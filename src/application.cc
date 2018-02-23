#include "application.h"

namespace arche {

void Application::handleHelp(const std::string &name,
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

int Application::main(const std::vector<std::string> &args) {
  if (!_helpRequested) {
    unsigned short port =
        (unsigned short)config().getInt("Application.port", 8080);
    std::string format(config().getString("Application.format",
                                          DateTimeFormat::SORTABLE_FORMAT));

    ServerSocket svs(port);
    HTTPServer srv(new Router(format), svs, new HTTPServerParams);
    srv.start();
    waitForTerminationRequest();
    srv.stop();
  }
  return Application::EXIT_OK;
}

void Application::defineOptions(OptionSet &options) {
  ServerApplication::defineOptions(options);

  options.addOption(Option("help", "h", "display argument help information")
                        .required(false)
                        .repeatable(false)
                        .callback(OptionCallback<Application>(
                            this, &Application::handleHelp)));
}

} // namespace arche
