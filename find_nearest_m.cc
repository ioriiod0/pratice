#include <iterator>
#include <map>
#include <set>
#include <utility>

template <typename BST,typename Key>
std::pair<BST::iterator,BST::iterator> find_nearest_m(const BST& bst,const Key& key,int m)
{
	if(bst.size() <= m)
		return std::make_pair(bst.begin(),bst.end());

	BST::iterator r = bst.lower_bound(key); //KEY右面的游标
	BST::iterator l = std::advance(r,-1); //KEY左面的游标
	
	while( m > 0 && r != bst.end() && l != bst.rend() )
	{
		if(l->first - key < key - r->first)//右面比较靠近
		{
			++r;
			m-=1;
		}
		else if( it->first - key > key - j->first) //左面比较近
		{
			--l;
			m-=1;
		}
		else //两边一样近
		{
			++r;
			--l;
			m-=2;
		}
 
	}

	if(m <= 0) //m可能为-1 ,因为31行
	{
		return std::make_pair(std::advance(l++,-m),r)
	}
	else if(l != bst.rend()) //m不为0，且左面不空，
	{
		return std::make_pair(std::advance(l++,-m),r);
	}
	else //m不为0，且右面不为空。
	{
		return std::make_pair(i++,std::advance(r,m));
	}

} 


int int main(int argc, char const *argv[])
{
	

	return 0;
}