#include <iostream>
#include <iomanip>
#include <omp.h>

using namespace std;

double calcPi(long long numsteps) {
  double x, pi = 0;
  
  for (int i = 0; i < numsteps; ++i) {
    x = i*1. / numsteps;
    pi += 4. / (1. + x*x);
  }
  return pi / numsteps;
}

double calcPiParalelo(long long numsteps) {
  double x, pi = 0;
  
  for (int i = 0; i < numsteps; ++i) {
    x = i*1. / numsteps;
    pi += 4. / (1. + x*x);
  }
  return pi / numsteps;
}

int main(int argc, char** argv) {
  if (argc < 2) return -1;
  
  double ini, fin;
  
  ini = omp_get_wtime();
  cout << "Pi = " << setprecision(10) << calcPi(atoi(argv[1]));
  fin = omp_get_wtime();
  cout << "\tTiempo: " << (fin-ini) << endl;
  
  ini = omp_get_wtime();
  cout << "Pi = " << setprecision(10) << calcPiParalelo(atoi(argv[1]));
  fin = omp_get_wtime();
  cout << "\tTiempo: " << (fin-ini) << endl;
  
	return 0;
}
