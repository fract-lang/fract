#ifndef __CODE_FILE_HH
#define __CODE_FILE_HH

#include <fstream>
#include <iostream>
#include <vector>

#include "code_line.hh"

namespace Fract::Objects {
/// @brief Code file instance.
class code_file {
public:
  /// @brief File stream.
  std::ifstream stream;

  /// @brief Path of file.
  std::string path;

  /// @brief Code lines of file.
  std::vector<code_line> lines;
};
}

#endif // __CODE_FILE_HH
