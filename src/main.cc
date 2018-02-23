#include "server.h"

int main(int argc, char **argv) {
  arche::HTTPTimeServer app;
  return app.run(argc, argv);
}
