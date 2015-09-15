package com.cardinfolink.yunshouyin.util;

import java.util.Comparator;

import org.apache.http.NameValuePair;

public class ComparatorNameValuePair implements Comparator {

	@Override
	public int compare(Object lhs, Object rhs) {
		NameValuePair nameValuePair1=(NameValuePair) lhs;
		NameValuePair nameValuePair2=(NameValuePair) rhs;
		return nameValuePair1.getName().compareTo(nameValuePair2.getName());
	}

}
