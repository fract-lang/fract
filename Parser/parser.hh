#ifndef __PARSER_HH
#define __PARSER_HH

#include <iostream>
#include <vector>

#include "../Objects/code_line.hh"

using namespace Fract::Objects;

namespace Fract::Parser {

#define type_comment 999
#define type_function 1000
#define type_equals 1001
#define type_let 1002
#define type_name 1003
#define type_dotted_name 1004
#define type_value_setter 1005
#define type_value 1006
#define type_data_type 1007
#define type_end_type 1008
#define type_return 1009
#define type_import 1010
#define type_std_import 1011
#define type_if 1012
#define type_else_if 1013
#define type_else 1014
#define type_for 1015
#define type_while 1016
#define type_delete 1017
#define type_int16 1018
#define type_int32 1019
#define type_int64 1020
#define type_int64_64 1021
#define type_unsigned_int16 1022
#define type_unsigned_int32 1023
#define type_unsigned_int64 1024
#define type_unsigned_int64_64 1025
#define type_float 1026
#define type_double 1027
#define type_boolean 1028
#define type_byte 1029
#define type_signed_byte 1030


/// @brief Parser of Fract.
class parser {
public:
  /**
   * @brief Ready lines to process.
   * @param lines Lines to ready.
   * @returns Ready lines.
  */
  static std::vector<code_line> readyLines(std::vector<std::string> lines);
};
} // namespace Fract::Parser

#endif // __PARSER_HH
