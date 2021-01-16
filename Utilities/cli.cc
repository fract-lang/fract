#include "cli.hh"

using namespace Fract::Objects;
using namespace Fract::Utilities;

void cli::printMapAsTable(std::map<std::string, std::string> map) {
  int maxlen = 0;
  for (std::map<std::string, std::string>::iterator it = map.begin();
       it != map.end(); ++it) {
    maxlen = maxlen < it->first.length() ? it->first.length() : maxlen;
  }
  maxlen += 5;
  for (std::map<std::string, std::string>::iterator it = map.begin();
       it != map.end(); ++it) {
    std::cout << it->first << std::string(maxlen - it->first.length(), ' ')
              << it->second << std::endl;
  }
}
