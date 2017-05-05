#include <iostream>
#include <omp.h>

using namespace std;

int main() {
  omp_set_num_threads(4);
#pragma omp parallel
{  
  int ID = omp_get_thread_num();
  cout << "hello " << ID;
  cout << ", world " << ID << endl;
}

	return 0;
}
