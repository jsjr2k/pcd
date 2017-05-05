#include <iostream>
#include <iomanip>
#include <omp.h>

//
// Trabaje con unos 100 000 000 de iteraciones.
//

using namespace std;

double calcPi(long long numsteps) {
  double x, pi = 0;
  
  for (int i = 0; i < numsteps; ++i) {
    x = i*1. / numsteps;
    pi += 4. / (1. + x*x);
  }
  return pi / numsteps;
}

int main(int argc, char** argv) {
  if (argc < 2) return -1;
  
  cout << "Pi = " << setprecision(10) << calcPi(atoi(argv[1])) << endl;
	return 0;
}
