package demo.example.mamh.fragmentdemo;

import android.app.Activity;
import android.app.FragmentManager;
import android.app.FragmentTransaction;
import android.graphics.Color;
import android.os.Bundle;
import android.view.View;
import android.view.Window;
import android.widget.ImageView;
import android.widget.TextView;

/**
 * 项目的主Activity，所有的Fragment都嵌入在这里。
 * 这个类中的注释已经写得非常详细了，下面我再带大家简单梳理一遍。在onCreate()方法中先是调用了initViews()
 * 来获取每个控件的实例，并给相应的控件设置好点击事件，然后调用setTabSelection()方法设置默认的选中项，
 * 这里传入的0说明默认选中第1个Tab项。
 * <p/>
 * 那么setTabSelection()方法中又是如何处理的呢？可以看到，首先第一步是调用clearSelection()方法来清理
 * 掉之前的选中状态，然后开启一个Fragment事务，并隐藏掉所有的Fragment，以防止有多个Fragment显示在界面上。
 * 接下来根据传入的index参数判断出选中的是哪一个Tab项，并改变该Tab项的图标和文字颜色，然后将相应的Fragment添
 * 加到界面上。这里注意一个细节，我们添加Fragment的时候并没有使用replace()方法，而是会先判断一下该Fragment是
 * 否为空，如果是空的则调用add()方法添加一个进来，如果不是空的则直接调用show()方法显示出来即可。那么为什么没有
 * 使用replace()方法呢？这是因为replace()方法会将被替换掉的那个Fragment彻底地移除掉，该Fragment的生命周期就
 * 结束了。当再次点击刚才那个Tab项的时候，就会让该Fragment的生命周期重新开始，onCreate()、onCreateView()等
 * 方法都会重新执行一遍。这显然不是我们想要的，也和ActivityGroup的工作原理不符，因此最好的解决方案就是使
 * 用hide()和show()方法来隐藏和显示Fragment，这就不会让Fragment的生命周期重走一遍了。
 * <p/>
 * 设置完默认选中项后，我们当然还可以通过点击Tab项来自由地切换界面，这就会进入到onClick()方法中。onClick()方
 * 法中的逻辑判断非常简单，当点击了消息标签时就会选中第1个tab项，点击联系人标签时就会选中第2个tab项，点击动态
 * 标签时就会选中第3个tab项，点击设置标签时就会选中第4个tab项。都是通过调用setTabSelection()方法来完成的，
 * 只是传入了不同的参数。
 */
public class MainActivity extends Activity implements View.OnClickListener {

    /**
     * 用于展示消息的Fragment
     */
    private MessageFragment messageFragment;

    /**
     * 用于展示联系人的Fragment
     */
    private ContactsFragment contactsFragment;

    /**
     * 用于展示动态的Fragment
     */
    private NewsFragment newsFragment;

    /**
     * 用于展示设置的Fragment
     */
    private SettingFragment settingFragment;

    /**
     * 消息界面布局
     */
    private View messageLayout;

    /**
     * 联系人界面布局
     */
    private View contactsLayout;

    /**
     * 动态界面布局
     */
    private View newsLayout;

    /**
     * 设置界面布局
     */
    private View settingLayout;

    /**
     * 在Tab布局上显示消息图标的控件
     */
    private ImageView messageImage;

    /**
     * 在Tab布局上显示联系人图标的控件
     */
    private ImageView contactsImage;

    /**
     * 在Tab布局上显示动态图标的控件
     */
    private ImageView newsImage;

    /**
     * 在Tab布局上显示设置图标的控件
     */
    private ImageView settingImage;

    /**
     * 在Tab布局上显示消息标题的控件
     */
    private TextView messageText;

    /**
     * 在Tab布局上显示联系人标题的控件
     */
    private TextView contactsText;

    /**
     * 在Tab布局上显示动态标题的控件
     */
    private TextView newsText;

    /**
     * 在Tab布局上显示设置标题的控件
     */
    private TextView settingText;

    /**
     * 用于对Fragment进行管理
     */
    private FragmentManager fragmentManager;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        requestWindowFeature(Window.FEATURE_NO_TITLE);
        setContentView(R.layout.activity_main);
        // 初始化布局元素
        initViews();
        fragmentManager = getFragmentManager();
        // 第一次启动时选中第0个tab
        setTabSelection(0);
    }

    /**
     * 在这里获取到每个需要用到的控件的实例，并给它们设置好必要的点击事件。
     */
    private void initViews() {
        messageLayout = findViewById(R.id.message_layout);
        contactsLayout = findViewById(R.id.contacts_layout);
        newsLayout = findViewById(R.id.news_layout);
        settingLayout = findViewById(R.id.setting_layout);
        messageImage = (ImageView) findViewById(R.id.message_image);
        contactsImage = (ImageView) findViewById(R.id.contacts_image);
        newsImage = (ImageView) findViewById(R.id.news_image);
        settingImage = (ImageView) findViewById(R.id.setting_image);
        messageText = (TextView) findViewById(R.id.message_text);
        contactsText = (TextView) findViewById(R.id.contacts_text);
        newsText = (TextView) findViewById(R.id.news_text);
        settingText = (TextView) findViewById(R.id.setting_text);
        messageLayout.setOnClickListener(this);
        contactsLayout.setOnClickListener(this);
        newsLayout.setOnClickListener(this);
        settingLayout.setOnClickListener(this);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.message_layout:
                // 当点击了消息tab时，选中第1个tab
                setTabSelection(0);
                break;
            case R.id.contacts_layout:
                // 当点击了联系人tab时，选中第2个tab
                setTabSelection(1);
                break;
            case R.id.news_layout:
                // 当点击了动态tab时，选中第3个tab
                setTabSelection(2);
                break;
            case R.id.setting_layout:
                // 当点击了设置tab时，选中第4个tab
                setTabSelection(3);
                break;
            default:
                break;
        }
    }

    /**
     * 根据传入的index参数来设置选中的tab页。
     *
     * @param index 每个tab页对应的下标。0表示消息，1表示联系人，2表示动态，3表示设置。
     */
    private void setTabSelection(int index) {
        // 每次选中之前先清楚掉上次的选中状态
        clearSelection();
        // 开启一个Fragment事务
        FragmentTransaction transaction = fragmentManager.beginTransaction();
        // 先隐藏掉所有的Fragment，以防止有多个Fragment显示在界面上的情况
        hideFragments(transaction);
        switch (index) {
            case 0:
                // 当点击了消息tab时，改变控件的图片和文字颜色
                messageImage.setImageResource(R.drawable.message_selected);
                messageText.setTextColor(Color.WHITE);
                if (messageFragment == null) {
                    // 如果MessageFragment为空，则创建一个并添加到界面上
                    messageFragment = new MessageFragment();
                    transaction.add(R.id.content, messageFragment);
                } else {
                    // 如果MessageFragment不为空，则直接将它显示出来
                    transaction.show(messageFragment);
                }
                break;
            case 1:
                // 当点击了联系人tab时，改变控件的图片和文字颜色
                contactsImage.setImageResource(R.drawable.contacts_selected);
                contactsText.setTextColor(Color.WHITE);
                if (contactsFragment == null) {
                    // 如果ContactsFragment为空，则创建一个并添加到界面上
                    contactsFragment = new ContactsFragment();
                    transaction.add(R.id.content, contactsFragment);
                } else {
                    // 如果ContactsFragment不为空，则直接将它显示出来
                    transaction.show(contactsFragment);
                }
                break;
            case 2:
                // 当点击了动态tab时，改变控件的图片和文字颜色
                newsImage.setImageResource(R.drawable.news_selected);
                newsText.setTextColor(Color.WHITE);
                if (newsFragment == null) {
                    // 如果NewsFragment为空，则创建一个并添加到界面上
                    newsFragment = new NewsFragment();
                    transaction.add(R.id.content, newsFragment);
                } else {
                    // 如果NewsFragment不为空，则直接将它显示出来
                    transaction.show(newsFragment);
                }
                break;
            case 3:
            default:
                // 当点击了设置tab时，改变控件的图片和文字颜色
                settingImage.setImageResource(R.drawable.setting_selected);
                settingText.setTextColor(Color.WHITE);
                if (settingFragment == null) {
                    // 如果SettingFragment为空，则创建一个并添加到界面上
                    settingFragment = new SettingFragment();
                    transaction.add(R.id.content, settingFragment);
                } else {
                    // 如果SettingFragment不为空，则直接将它显示出来
                    transaction.show(settingFragment);
                }
                break;
        }
        transaction.commit();
    }

    /**
     * 清除掉所有的选中状态。
     */
    private void clearSelection() {
        messageImage.setImageResource(R.drawable.message_unselected);
        messageText.setTextColor(Color.parseColor("#82858b"));
        contactsImage.setImageResource(R.drawable.contacts_unselected);
        contactsText.setTextColor(Color.parseColor("#82858b"));
        newsImage.setImageResource(R.drawable.news_unselected);
        newsText.setTextColor(Color.parseColor("#82858b"));
        settingImage.setImageResource(R.drawable.setting_unselected);
        settingText.setTextColor(Color.parseColor("#82858b"));
    }

    /**
     * 将所有的Fragment都置为隐藏状态。
     *
     * @param transaction 用于对Fragment执行操作的事务
     */
    private void hideFragments(FragmentTransaction transaction) {
        if (messageFragment != null) {
            transaction.hide(messageFragment);
        }
        if (contactsFragment != null) {
            transaction.hide(contactsFragment);
        }
        if (newsFragment != null) {
            transaction.hide(newsFragment);
        }
        if (settingFragment != null) {
            transaction.hide(settingFragment);
        }
    }
}
