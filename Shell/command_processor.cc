#include "command_processor.hh"

using namespace Fract::Shell;
using namespace Fract::Utilities;

std::string
command_processor::getNamespace(std::string cmd) {
  std::size_t pos = cmd.find(" ");
  return pos == std::string::npos ? cmd : cmd.substr(0, pos);
}

std::string
command_processor::removeNamespace(std::string cmd) {
  std::size_t pos = cmd.find(" ");
  return pos == std::string::npos ? "" : cmd.substr(pos + 1);
}

bool
command_processor::getArguments(std::string cmd,
                                  std::vector<std::string>* dest) {
  std::smatch match;
  while (
      std::regex_search(cmd, match, std::regex("(^|\\s+)-\\w+(?=($|\\s+))"))) {
    std::string arg = string::toLower(string::trim(match[0]));
    if (std::find(dest->begin(), dest->end(), arg) != dest->end()) {
      std::cout << "A argument cannot be written more than once!" << std::endl;
      return false;
    }
    dest->push_back(arg);
    cmd = match.suffix();
  }
  return true;
}

std::string
command_processor::removeArguments(std::string cmd) {
  return std::regex_replace(cmd, std::regex("(^|\\s+)-\\w+(?=($|\\s+))"), "");
}
