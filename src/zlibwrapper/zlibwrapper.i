%module zlibwrapper

%include <typemaps.i>
%include "std_vector.i"

namespace std {
   %template(BinaryData) vector<unsigned char>;
}

//%{
/* Include in the generated wrapper file */
//#include "zlibwrapper.hpp"
//%}

/* Tell SWIG about it */
//%include "zlibwrapper.hpp"


%inline %{
extern std::vector< unsigned char > uncompress(const std::vector< unsigned char > &sourceData);
extern std::vector< unsigned char > compress(const std::vector< unsigned char > &sourceData);
%}

extern std::vector< unsigned char > uncompress(const std::vector< unsigned char > &sourceData);
extern std::vector< unsigned char > compress(const std::vector< unsigned char > &sourceData);
