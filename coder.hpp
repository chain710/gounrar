/****************************************************************************
 *  Contents: 'Carryless rangecoder' by Dmitry Subbotin                     *
 ****************************************************************************/
#pragma once
#include "rartypes.hpp"

class Unpack;
class RangeCoder
{
  public:
    void InitDecoder(Unpack *UnpackRead);
    int GetCurrentCount();
    uint GetCurrentShiftCount(uint SHIFT);
    void Decode();
    void PutChar(unsigned int c);
    unsigned int GetChar();

    uint low, code, range;
    struct SUBRANGE 
    {
      uint LowCount, HighCount, scale;
    } SubRange;

    Unpack *UnpackRead;
};


// (int) cast before "low" added only to suppress compiler warnings.
#define ARI_DEC_NORMALIZE(code,low,range,read)                           \
{                                                                        \
  while ((low^(low+range))<TOP || range<BOT && ((range=-(int)low&(BOT-1)),1)) \
  {                                                                      \
    code=(code << 8) | read->GetChar();                                  \
    range <<= 8;                                                         \
    low <<= 8;                                                           \
  }                                                                      \
}
