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
  static std::string boldRed;
  static std::string boldGreen;
  static std::string boldYellow;
  static std::string boldMagenta;
  static std::string boldCyan;
  static std::string reset;
  static byte min;
  static byte max;

  /// @brief Enable virtual terminal processing.
  static void enableVTP();

  /**
   * @brief Create ANSI color code by color.
   * @param color Color instance.
   * @return ANSI code of rgb values.
   */
  static std::string toANSI(color color);

  /**
   * @brief Create ANSI color code by rgb.
   * @param r Red.
   * @param g Green.
   * @param b Blue.
   * @return ANSI code of rgb values.
   */
  static std::string toANSI(byte r, byte g, byte b);

  /// @brief Red.
  byte r = max;
  /// @brief Green.
  byte g = max;
  /// @brief Blue.
  byte b = max;

  /**
   * @brief Create new instance.
   * @param r Red.
   * @param g Green.
   * @param b Blue.
   */
  color(byte r, byte g, byte b);
};
}  // namespace Fract::Objects

#endif  // __COLOR_HH
