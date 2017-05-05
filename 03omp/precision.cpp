#include <iostream>
#include <iomanip>

using namespace std;

int main() {
  double a = 1000000000.;
  double b = 0.000001;
  double c = a;
  for (int i = 0; i < 1000000; ++i) {
    c += b;
  }
  c -= a;

  cout << c << endl;

  return 0;
}
