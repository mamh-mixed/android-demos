package com.cardinfo.framelib.model;


import java.util.List;
import org.apache.http.NameValuePair;
public class RequestParam {
  private String url;
  private  List<NameValuePair> params;
public String getUrl() {
	return url;
}
public void setUrl(String url) {
	this.url = url;
}
public List<NameValuePair> getParams() {
	return params;
}
public void setParams(List<NameValuePair> params) {
	this.params = params;
}

  
}
