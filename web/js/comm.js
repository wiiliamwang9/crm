const {createApp, ref, onMounted} = Vue
const app = createApp({
    setup() {
        const userInfo = ref()
        const todayIndex = ref(0)
        const todayIndexArray = ref(['All', 'Todo', 'Regular', 'Sampled', 'Shipped'])
        const todayList = ref([])
        const todayTotal = ref(0)
        const upcomingIndex = ref(0)
        const upcomingIndexArray = ref(['All', 'Stale', 'Never'])
        const upcomingList = ref([])
        const upcomingTotal = ref(0)
        const userId = ref(1)
        const tagColor = ref(['big-customer', 'new-customer', 'never-ordered'])
        onMounted(() => {
            // 等待配置加载完成
            const waitForConfig = () => {
                if (window.configLoaded && window.baseUrl) {
                    getTodayData("All")
                    getUpcomingData("All")
                    getUserInfo(userId.value)
                } else {
                    setTimeout(waitForConfig, 10);
                }
            };
            waitForConfig();
        })
        //今日待跟进的数据
        const getTodayData = (statusFilter) => {
            if (!window.baseUrl) {
                console.error('window.baseUrl is not defined');
                return;
            }
            axios.post(window.baseUrl + "/api/v1/dashboard/search", {
                user_id: userId.value,
                time_filter: "Today",
                status_filter: statusFilter,
                show_all: false,
                page_size: 4,
                page: 1
            }).then(res => {
                if (res.status == 200) {
                    console.log(res.data.data.list)
                    todayList.value = res.data.data.list
                    todayTotal.value = res.data.data.total
                } else {

                }
            })
        }
        //近期
        const getUpcomingData = (statusFilter) => {
            if (!window.baseUrl) {
                console.error('window.baseUrl is not defined');
                return;
            }
            axios.post(window.baseUrl + "/api/v1/dashboard/search", {
                user_id: userId.value,
                time_filter: "Upcoming",
                status_filter: statusFilter,
                show_all: false,
                page_size: 4,
                page: 1
            }).then(res => {
                if (res.status == 200) {
                    console.log(res.data.data.list)
                    upcomingList.value = res.data.data.list
                    upcomingTotal.value = res.data.data.total
                } else {

                }
            })
        }
        //用户
        const getUserInfo = (id) => {
            if (!window.baseUrl) {
                console.error('window.baseUrl is not defined');
                return;
            }
            axios.get(window.baseUrl + "/api/v1/users/" + id).then(res => {
                if (res.status == 200) {
                    userInfo.value = res.data.data
                } else {

                }
            })
        }
        //tabs切换
        const onTodayChange = (e) => {
            todayIndex.value = e
            getTodayData(todayIndexArray.value[e])
        }
        const onUpcomingChange = (e) => {
            upcomingIndex.value = e
            getUpcomingData(upcomingIndexArray.value[e])
        }
        //随机数
        const getRandomIntInclusive = (min, max) => {
            min = Math.ceil(min);
            max = Math.floor(max);
            return Math.floor(Math.random() * (max - min + 1)) + min;
        }
        //详情
        const onRead = (id) => {
            window.location.href = './customer.html?id=' + id
        }
        
        // 弹窗相关
        const showActionMenu = ref(false)
        const showTodoDialog = ref(false)
        const showSOSDialog = ref(false)
        const currentCustomerId = ref(null)
        
        // 点击Frame.png显示操作菜单
        const onMoreClick = (customerId) => {
            currentCustomerId.value = customerId
            showActionMenu.value = true
        }
        
        // 操作菜单选项点击
        const onAddTodo = () => {
            showActionMenu.value = false
            showTodoDialog.value = true
        }
        
        const onSOS = () => {
            showActionMenu.value = false
            showSOSDialog.value = true
        }
        
        const onCall = () => {
            showActionMenu.value = false
            // 这里可以添加打电话的逻辑
            alert('拨打电话功能')
        }
        
        return {
            todayList,
            getRandomIntInclusive,
            todayTotal,
            tagColor,
            todayIndex,
            onTodayChange,
            upcomingList,
            upcomingTotal,
            upcomingIndex,
            onUpcomingChange,
            userInfo,
            onRead,
            showActionMenu,
            showTodoDialog,
            showSOSDialog,
            currentCustomerId,
            onMoreClick,
            onAddTodo,
            onSOS,
            onCall
        }
    }
})
app.mount('#app')

const app1 = createApp({
    setup() {
        const customerInfo = ref()
        const customerStats = ref()
        const tagColor = ref(['big-customer', 'new-customer', 'never-ordered'])
        const show1 = ref(false)
        const show2 = ref(false)
        const show3 = ref(false)
        const show4 = ref(false)
        const show5 = ref(false)
        const form1 = ref({
            time: "明日",
            type: '日常跟进',
            content: "",
            switch: false,
            tip: '企微',
            peper: '我自己'
        })
        const form2 = ref({
            time: "明日",
            type: '电话跟进',
            content: "",
            switch: false
        })
        const form3 = ref({
            shuren: "徐晓二",
            guanxi: "好友",
            content: ""
        })
        const form4 = ref({
            shuren: "徐晓二",
            guanxi: "好友",
            content: ""
        })
        const tabIndex = ref(0)
        onMounted(() => {
            const params = new URLSearchParams(new URL(window.location.href).search);
            const id = params.get('id');
            // 等待配置加载完成
            const waitForConfig = () => {
                if (window.configLoaded && window.baseUrl) {
                    getCustomerInfo(id)
                    getCustomerStats(id)
                    getTodosData(id)
                } else {
                    setTimeout(waitForConfig, 10);
                }
            };
            waitForConfig();
        })
        //客户详情
        const getCustomerInfo = (id) => {
            if (!window.baseUrl) {
                console.error('window.baseUrl is not defined');
                return;
            }
            axios.get(window.baseUrl + "/api/v1/customers/" + id).then(res => {
                console.log(res)
                if (res.status == 200) {
                    customerInfo.value = res.data.data
                } else {

                }
            })
        }
        //数据
        const getCustomerStats = (id) => {
            if (!window.baseUrl) {
                console.error('window.baseUrl is not defined');
                return;
            }
            axios.post(window.baseUrl + "/api/v1/customers/stats", {
                customer_id: parseInt(id)
            }).then(res => {
                console.log(res)
                if (res.status == 200) {
                    customerStats.value = res.data.data
                } else {

                }
            })
        }
        //待办
        const getTodosData = (id) => {
            if (!window.baseUrl) {
                console.error('window.baseUrl is not defined');
                return;
            }
            axios.get(window.baseUrl + "/api/v1/todos", {
                params: {
                    customer_id: id,
                }
            }).then(res => {
                console.log(res)
                // if(res.status == 200){
                //     customerInfo.value = res.data.data
                // }else{

                // }
            })
        }
        //tab切换
        const onTabChange = (index) => {
            tabIndex.value = index
        }
        //随机数
        const getRandomIntInclusive = (min, max) => {
            min = Math.ceil(min);
            max = Math.floor(max);
            return Math.floor(Math.random() * (max - min + 1)) + min;
        }
        return {
            show1,
            show2,
            show3,
            show4,
            show5,
            form1,
            form2,
            form3,
            form4,
            customerInfo,
            customerStats,
            tabIndex,
            onTabChange,
            getRandomIntInclusive,
            tagColor
        }
    }
})
app1.use(vant);
app1.mount('#app1')

const app2 = createApp({
    setup() {
        const userId = ref(1)
        const dataList = ref([])
        const dataTotal = ref(0)
        const tagColor = ref(['big-customer', 'new-customer', 'never-ordered'])
        const type1Index = ref(0)
        const type1IndexArray = ref(['All', 'Today', 'Upcoming'])
        const type2Index = ref(0)
        const type2IndexArray = ref(['All', 'Todo', 'Regular', 'Sampled', 'Shipped'])
        const show1 = ref(false)
        const form1 = ref({
            time: '',
            type: '',
            content: ''
        })
        onMounted(() => {
            console.log("ddd")
            // 等待配置加载完成
            const waitForConfig = () => {
                if (window.configLoaded && window.baseUrl) {
                    getData("All", "All")
                } else {
                    setTimeout(waitForConfig, 10);
                }
            };
            waitForConfig();
        })
        //改变
        const onType1Change = (e) => {
            type1Index.value = e
            getData(type1IndexArray.value[e], type2IndexArray.value[type2Index.value])
        }
        const onType2Change = (e) => {
            type2Index.value = e
            getData(type1IndexArray.value[type1Index.value], type2IndexArray.value[e])
        }
        //数据
        const getData = (timeFilter, statusFilter) => {
            if (!window.baseUrl) {
                console.error('window.baseUrl is not defined');
                return;
            }
            axios.post(window.baseUrl + "/api/v1/dashboard/search", {
                user_id: userId.value,
                time_filter: timeFilter,
                status_filter: statusFilter,
                show_all: false,
                page_size: 9999,
                page: 1
            }).then(res => {
                if (res.status == 200) {
                    console.log(res.data.data.list)
                    dataList.value = res.data.data.list
                    dataTotal.value = res.data.data.total
                } else {

                }
            })
        }
        //详情
        const onRead = (id) => {
            window.location.href = './customer.html?id=' + id
        }
        //随机数
        const getRandomIntInclusive = (min, max) => {
            min = Math.ceil(min);
            max = Math.floor(max);
            return Math.floor(Math.random() * (max - min + 1)) + min;
        }
        return {
            dataList,
            dataTotal,
            getRandomIntInclusive,
            tagColor,
            onType2Change,
            type2Index,
            type1Index,
            onType1Change,
            onRead,
            show1,
            form1
        }
    }
})
app2.use(vant);
app2.mount('#app2')