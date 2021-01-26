#ifndef __RMRX_RANGE_RESULT_HH
#define __RMRX_RANGE_RESULT_HH

#include <iostream>
#include <vector>

#include "../Objects/token.hh"

using namespace Fract::Objects;

namespace Fract::Parser {
/// @brief Result instance of formatter::rmlx_range
struct rmlx_range_result {
public:
  bool found; /* Range is found */
  std::vector<token> range; /* Tokens of range*/
  int index; /* Index of replace iterator of original vector */
};
} // namespace Fract::Parser

#endif // __RMRX_RANGE_RESULT_HH
