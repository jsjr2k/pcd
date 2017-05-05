#include <iostream>
#include <iomanip>
#include <omp.h>

#define MAXTH 8

using namespace std;

double calcPi(long long numsteps) {
  double x, pi = 0;
  double step = 1. /numsteps;

  for (int i = 0; i < numsteps; ++i) {
    x = i * step;
    pi += 4. / (1. + x*x);
  }
  return pi / numsteps;
}

double calcPiParalelo(long long numsteps) {
  omp_set_num_threads(MAXTH);

  double pi = 0;
  double step = 1. /numsteps;

#pragma omp parallel for reduction(+:pi)
  for (int i = 0; i < numsteps; ++i) {
    double x = i * step;
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
