#ifndef __COLOR_HH
#define __COLOR_HH

#include <iostream>

#if _WIN32
#include <windows.h>
#endif  // _WIN32

namespace Fract::Objects {
/// @brief CLI Color.
struct color {
  public:
  static std::string white;
  static std::string red;
  static std::string green;
  static std::string yellow;
  static std::string blue;
  static std::string magenta;
  static std::string cyan;
  static std::string bold_red;
  static std::string bold_green;
  static std::string bold_yellow;
  static std::string bold_magenta;
  static std::string bold_cyan;
  static std::string reset;

  /// @brief Enable virtual terminal processing.
  static void enable_vtp();

  /**
   * @brief Create ANSI color code by color.
   * @param color Color instance.
   * @return ANSI code of rgb values.
   */
  static std::string to_ansi(color color);

  /**
   * @brief Create ANSI color code by rgb.
   * @param r Red.
   * @param g Green.
   * @param b Blue.
   * @return ANSI code of rgb values.
   */
  static std::string to_ansi(unsigned short r, unsigned short g, unsigned short b);

  /// @brief Red.
  unsigned short r;
  /// @brief Green.
  unsigned short g;
  /// @brief Blue.
  unsigned short b;

  /**
   * @brief Create new instance.
   * @param r Red.
   * @param g Green.
   * @param b Blue.
   */
  color(unsigned short r, unsigned short g, unsigned short b);
};
}  // namespace Fract::Objects

#endif  // __COLOR_HH
