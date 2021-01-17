#include "file_system.hh"

using namespace Fract::Utilities;

char* file_system::_WORKING_DIR_ = getcwd(NULL, 0);

bool file_system::existFile(std::string path) {
  std::ifstream file(path);
  bool state = file.good();
  delete &file;
  return state;
}

std::vector<std::string> file_system::getLinesOfFile(std::string path) {
  std::ifstream file(path);
  if (!file.good()) {
    std::cout << "File is not found this path: " << path << std::endl;
    exit(1);
  }
  std::vector<std::string> lines;
  std::string line;
  while (std::getline(file, line)) lines.push_back(line);
  file.close();
  delete &file;
  return lines;
}