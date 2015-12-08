package com.example.mamh.edittextclear;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;


/**
 * 今天给大家带来一个很实用的小控件ClearEditText，就是在Android系统的输入框右边加入一个小图标，点击小图标可以
 * 清除输入框里面的内容，IOS上面直接设置某个属性就可以实现这一功能，但是Android原生EditText不具备此功能，所以
 * 要想实现这一功能我们需要重写EditText，接下来就带大家来实现这一小小的功能
 * 我们知道，我们可以为我们的输入框在上下左右设置图片，所以我们可以利用属性android:drawableRight设置我们的删除
 * 小图标，如图
 * 我这里设置了左边和右边的图片，如果我们能为右边的图片设置监听，点击右边的图片清除输入框的内容并隐藏删除图标，这样
 * 子这个小功能就迎刃而解了，可是Android并没有给允许我们给右边小图标加监听的功能，这时候你是不是发现这条路走不通呢，
 * 其实不是，我们可能模拟点击事件，用输入框的的onTouchEvent()方法来模拟，
 * <p/>
 * 当我们触摸抬起（就是ACTION_UP的时候）的范围  大于输入框左侧到清除图标左侧的距离，小与输入框左侧到清除图片右侧的
 * 距离，我们则认为是点击清除图片，当然我这里没有考虑竖直方向，只要给清除小图标就上了监听，其他的就都好处理了，我先把
 * 代码贴上来，在讲解下
 */
public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
    }
}
