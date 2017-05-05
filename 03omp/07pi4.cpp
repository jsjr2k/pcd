#include <iostream>
#include <iomanip>
#include <omp.h>

#define MAXTH 8

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
  omp_set_num_threads(MAXTH);

  double x, pi = 0;
  double pilocals[MAXTH];

#pragma omp parallel
{
  int id = omp_get_thread_num();
  double x;
  double step = 1. / numsteps;
  int bloque = numsteps / MAXTH;

  pilocals[id] = 0.;
  for (int i = id; i < numsteps; i += MAXTH) {
    x = i * step;
    pilocals[id] += 4. / (1. + x*x);
  }
}
  for (int i = 0; i < MAXTH; ++i) {
    pi += pilocals[i];
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
