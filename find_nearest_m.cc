/*
给一个N个节点的二叉搜索树（BST/Binary Search Tree），给一个Key，返回与key最接近的m个节点（m<N）。
*/

#include <iterator>
#include <map>
#include <iostream>
#include <utility>



template <typename K,typename V>
std::pair<typename std::map<K,V>::iterator,typename std::map<K,V>::iterator> 
find_nearest_n(std::map<K,V>& bst,const K& key,int m)
{
	typedef typename std::map<K,V>::iterator iterator;
	typedef typename std::map<K,V>::reverse_iterator reverse_iterator;

	if(bst.size() <= m)
		return std::make_pair(bst.begin(),bst.end());

	iterator r = bst.lower_bound(key); //KEY右面的游标
	reverse_iterator l(r); //KEY左面的游标
	reverse_iterator left_bound(bst.begin());//左边的边界
	
	while( m > 0 && r != bst.end() && l != left_bound)
	{
		if(r->first - key < key - l->first)//右面比较靠近
		{
			++r;
			m-=1;
		}
		else if(r->first - key > key - l->first) //左面比较近
		{
			++l;
			m-=1;
		}
		else //两边一样近,则预先选左边（也可以先选择右边）
		{
			++l;
			m-=1;
		}
 
	}

	if (m > 0)
	{
		if(l != left_bound) //m不为0，且左面不空，
		{
			while(l != left_bound && m-- > 0)
				++l;
		}
		else //m不为0，且右面不为空。
		{
			while(r != bst.end() && m-- > 0)
				++r;
		}

	}
	
	return std::make_pair(l.base(),r);

} 


int main(int argc, char const *argv[])
{
	
	typedef typename std::map<int,int>::iterator iterator;
	std::map<int,int> datas;


	datas[1] = 1;
	datas[3] = 3;
	datas[4] = 4;
	datas[5] = 5;
	datas[8] = 8;
	datas[10] = 10;


	{
		std::pair<iterator,iterator> ret = find_nearest_n(datas,0,0);

		for(iterator it = ret.first;it != ret.second;++it)
		{
			std::cout<<it->first<<std::endl;
		}
		std::cout<<"===========NONE============="<<std::endl;

	}

	{
		std::pair<iterator,iterator> ret = find_nearest_n(datas,6,3);

		for(iterator it = ret.first;it != ret.second;++it)
		{
			std::cout<<it->first<<std::endl;
		}
		std::cout<<"========================"<<std::endl;

	}

	{
		std::pair<iterator,iterator> ret = find_nearest_n(datas,5,3);

		for(iterator it = ret.first;it != ret.second;++it)
		{
			std::cout<<it->first<<std::endl;
		}
		std::cout<<"========================"<<std::endl;

	}


	{
		std::pair<iterator,iterator> ret = find_nearest_n(datas,11,3);

		for(iterator it = ret.first;it != ret.second;++it)
		{
			std::cout<<it->first<<std::endl;
		}
		std::cout<<"========================"<<std::endl;
	}

	{

		std::pair<iterator,iterator> ret = find_nearest_n(datas,0,3);

		for(iterator it = ret.first;it != ret.second;++it)
		{
			std::cout<<it->first<<std::endl;
		}
		std::cout<<"========================"<<std::endl;

	}

	{
		
		std::pair<iterator,iterator> ret = find_nearest_n(datas,2,11);

		for(iterator it = ret.first;it != ret.second;++it)
		{
			std::cout<<it->first<<std::endl;
		}
		std::cout<<"========================"<<std::endl;

	}



	return 0;
}