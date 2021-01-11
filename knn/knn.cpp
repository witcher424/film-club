#include <vector>
#include <iostream>
#include "math.h"

using namespace::std;

//функция для среднего значения чтоб не засорять другие функции
vector<double> average(vector<vector<double>> a) {
    vector<double> sum;
    int cols=a.size();//юзеры
    int rows = a[0].size();//фильмы

    sum.resize(rows);

    for (size_t i = 0; i < rows; i++) {
        for (size_t j = 0; j < cols; j++) {
            sum[i] += a[j][i];
        }
        sum[i] = sum[i] / cols;
    }
    return sum;
}
//функция для срзнач но для числа это так тупо если честно извините за быдлокод я быдло
double aver_onedir(vector<double> a) {
    double sum=0;
    int rows = a.size();//фильмы

    for (size_t i = 0; i < rows; i++) {
            sum += a[i];
        }
        sum = sum/ rows;
    return sum;
}
//сумма произведений
double summult(vector<double> a, vector<double> b) {
    size_t n=b.size(), sum;
    sum = 0;
    for (size_t i = 0; i < n; i++) {
        sum = sum + a[i] * b[i];
    }
    return sum;
}

double summult2(double a, vector<double> b) {
    size_t n = b.size();
    double sum = 0;
    for (size_t i = 0; i < n; i++) {
        sum = sum + a * b[i];
    }
    return sum;
}

//просто сумма
double sum(vector<double> a) {
    size_t n = a.size();
    double sum = 0;
    for (size_t i = 0; i < n; i++) {
        sum = sum + a[i];
    }
    return sum;
}



//функция реализующая knn ищет соседей по косинусу, предсказывает рейтинг с поvощью GroupLens Algorithm
double knn(vector<double> ri, vector<vector<double>> rj) {
    vector<double> w;
    vector<double> pred_ri;
    w.resize(rj.size());
    pred_ri.resize(rj.size());
    cout << rj.size();


    for (size_t i = 0; i < w.size(); i++) {
        w[i] = summult(ri, rj[i]) / (sqrt(summult(ri, ri)) * sqrt(summult(rj[i], rj[i])));
        if (w[i] < 0.7) {//тут короче удаляю всех неблизких юзеров


            /*w.erase(w.begin() + j);
            for (size_t i = 0; i < rj.size(); i++) {
            rj[i].erase(rj[i].begin() + j);
            }*/
        }
    }


//дебагу нет конца


    pred_ri.resize(ri.size());

    for (size_t i = 0; i < pred_ri.size(); i++) {
        pred_ri[i] = 0;
        pred_ri[i] = average(rj)[i]+((summult2(w,rj)-sum(w)*aver_onedir(rj))/(sqrt(w*w)));}

    cout <<"predicted  "<< pred_ri[0]<<" "<< pred_ri[1]<<" "<<pred_ri[2];
    return 0;
}


int main()
{

    vector<double> urate;
    vector <vector<double>> rateothers;

    int film = 3, user = 2;
    //Grow rows by m
    rateothers.resize(user);
    for (int i = 0; i < user; ++i)
    {
        rateothers[i].resize(film);
    }

    rateothers[0][0] = 5;
    rateothers[0][1] = 3;
    rateothers[0][2] = 1;
    rateothers[1][0] = 4;
    rateothers[1][1] = 2;
    rateothers[1][2] = 2;
    //прога не понимает что делать при исключении юзеров исправить
 
    urate.resize(film);
    urate[0] = 1;
    urate[1] = 5;

    knn(urate, rateothers);
}
