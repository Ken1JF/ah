package ah_test

import (
	"fmt"
	. "gitHub.com/Ken1JF/ahgo/ah"
)

func ExamplePrintZKeys() {
	fmt.Println("keyCount =", KeyCount)
	fmt.Println("ZSeed =", ZSeed)
	var i, j uint8
	fmt.Println("BlackZKey =")
	for i = 0; i < MaxBoardSize; i += 1 {
		for j = 0; j < MaxBoardSize; j += 1 {
			fmt.Print(BlackZKey[i][j], ",")
		}
		fmt.Println()
	}
	fmt.Println("WhiteZKey =")
	for i = 0; i < MaxBoardSize; i += 1 {
		for j = 0; j < MaxBoardSize; j += 1 {
			fmt.Print(WhiteZKey[i][j], ",")
		}
		fmt.Println()
	}
	fmt.Println("KoZKey =")
	for i = 0; i < MaxBoardSize; i += 1 {
		for j = 0; j < MaxBoardSize; j += 1 {
			fmt.Print(KoZKey[i][j], ",")
		}
		fmt.Println()
	}
	//	fmt.Println("BlackToPlayZKey =", BlackToPlayZKey)
	fmt.Println("WhiteToPlayZKey =", WhiteToPlayZKey)
	// Output:
	// keyCount = 1084
	// ZSeed = 2309
	// BlackZKey =
	// 3518119828,2558606607,4028773274,365002044,3309554128,672890050,388554656,593403309,2608910757,3619422039,2221508301,3591828936,1508772138,951096920,3288463493,495429694,705902650,116625218,3354977286,
	// 3197386693,3072465261,2382795726,81904148,565996831,3666271361,1357359368,4060654627,137317671,1379054729,3379188104,3603685308,3508585185,3114106401,2060431485,689911161,404484884,1809156144,2965255689,
	// 236898027,3422734164,3562903709,2485265014,3044860503,3707919404,3631264180,2892765353,164524491,252958979,2273074205,315864474,4009202482,2799794835,2671134850,218925278,2694506510,1933922810,4180336915,
	// 2154424282,1167021653,2260793018,4004750091,3788452237,1046201883,1156138676,4006655988,90309034,9316532,3680940460,3381217062,2140597901,4060658091,3402741124,3761255313,1175623945,2305435233,3524348894,
	// 305528857,23498852,857133410,2142708676,1772936580,1264029056,1949764330,1053685526,3447025902,2288395870,2033408494,233360509,1140615687,2315413907,4153961436,344454384,2211053150,2440169704,522187649,
	// 1038417310,2430283863,1764865001,214539570,3698822063,3231192442,365192750,483553521,2230311860,2974696255,3593602878,3256341054,1898092466,1447526483,3743505621,887782763,2965453873,680870333,4043436775,
	// 1793871647,316454132,3473549662,2157158554,22792589,2381199185,4193077016,2652368939,3555340702,1612075935,1102595192,963201476,2533592065,711230539,1414483236,741842209,573537190,844048701,1997049508,
	// 2705682466,2571147551,94615564,880051341,564110602,4167271818,2556407757,3140924367,459702,2057423177,414595081,4022064521,316434806,3130508585,236391144,1361087629,3236786258,1053764789,2572108209,
	// 433760691,1307958989,457868606,3360605182,3027299634,2307633112,3406053403,2318447724,42727545,62257678,1161846808,3625464864,1912159924,3494194155,1539581916,970983835,1047864101,500912442,4068462897,
	// 495478141,725977348,1923920402,180900141,3394466886,725169215,168556538,93267850,491669059,1920337330,779336083,183364332,3389281391,3890758978,4243865107,1630937552,3802851104,1887887228,1114636679,
	// 2729002889,2076186515,943380618,3088004251,248496028,104465657,374487778,558458438,764628605,3172173801,572460684,1818893909,2486326313,1845166838,1101920237,2450792113,3179462723,3867626803,2581559176,
	// 1781647972,1560771890,3541420551,1734824377,27192245,3678371144,2727368464,3461529231,1017354282,3670406878,3259776205,4120835308,3111077397,2027755064,1507307030,1766335701,1258350529,2030487121,2918706099,
	// 2035796201,930738696,3062569280,1463795305,3298032880,3472103096,3214415644,1147735320,2069860748,1082644899,1271145182,2990507982,776435170,2932623997,1528937987,3756695210,3712458639,2142108941,2316902236,
	// 3418437012,3159194990,1115903297,3672684873,3051845649,46298960,4227767493,219700694,2289216453,4062939707,1666531386,3883011040,1597307401,3156478527,3984868163,640682385,3929523149,3809520585,3855865482,
	// 4294582213,1558616256,990027025,890940394,3565794357,3982817153,3290191297,1332220889,3201639266,4242718442,1166962654,814237143,960729880,3081909244,1680467131,2617392447,548271373,1155029872,3283875789,
	// 2358610799,2964720073,518054721,3764635664,87707438,613696364,252860879,1443789635,4056747915,2304015335,685919650,3492407233,2469979326,3451602321,2923885817,1345218781,2993978103,4014502921,4176643885,
	// 1503071521,1424349545,3418051765,3784859861,1157994162,1981061563,3263559760,1702716934,2057436917,905553179,3165289298,1339602642,1421698385,1072313488,3512791171,2572294651,4211496095,4158507577,3310021519,
	// 1720720639,1110262789,3403824915,1417455941,2680645438,1109278975,480638744,435149457,3304229506,1037367160,80448258,293930048,3057382486,239203202,1709626146,1486697838,4048062681,1245111806,2951505822,
	// 2104336513,4075609969,3324733800,1329310711,1175574308,673192966,3006581062,3562705613,4056790515,2222457763,1153843840,1955105759,3062351121,359958223,2448011047,713628542,1839260492,2518832543,141753109,
	// WhiteZKey =
	// 2922168950,2360010938,2331115303,3368447741,1447207927,14664226,2829336966,188586550,2585314721,2453408305,3271449857,723264217,3474480021,1272778768,3181004898,35796889,3135045433,3221578544,3206051120,
	// 687464366,182852028,2237413917,2742570375,2038655623,2758472103,905345839,1893976855,1151392013,2956787306,2641644005,4279846134,3849117922,140168,1866854649,2849531639,1426808917,2449494690,1614645345,
	// 2883282983,1954923899,2474019867,1827881142,1858347069,2116348437,450745708,2190107911,3037609413,1412027921,469371373,1889532220,1771469921,531228345,1512986354,3737379909,1767026351,1919502461,3085318272,
	// 676564621,950872500,2367640815,2943108880,3720051392,2921841148,1574734240,968428901,2283275944,3412207508,952790869,1168318421,349332303,4126643795,1006630985,2657967932,577937010,353958183,2739820085,
	// 943032642,1824237178,1939008574,1191196337,1276068244,3734823910,1719369117,1969602501,2916197876,3058062762,1945903297,3715828503,1756901408,1292644083,4279994560,1295822862,3057505873,1655719113,3677321169,
	// 1882547113,2715853588,4001496000,2213542342,921625298,3450473396,2927077121,3273507023,2114135809,549183324,45151999,2158037321,3997713928,1722780707,2661602662,2596896590,3959362332,3607186530,682451750,
	// 138061866,3281642198,3610107395,4143315092,2197829210,3813339223,1566267510,3659618828,4211008988,2930238881,888237712,3722788396,2393254092,1745267230,820306741,644420693,2356634332,697653628,4100132908,
	// 2891650243,3265460259,2932338834,3188964665,2961061501,475260662,1383601165,169365243,3407231701,4270835570,783647971,3956824761,2671209123,267138128,2159652668,3653403150,4209338790,3442289469,1077605323,
	// 1465763195,1429144617,3635331587,1800428884,1575672218,1805868067,3564024191,771288519,1628524224,4194137199,1492129636,2819523819,3047944370,1420317077,2892349975,3300500663,2581379090,1541395095,3021215791,
	// 833086738,3055757048,963539716,593149504,100427822,907815151,269012809,3114485370,471850487,3242274448,1946797578,3137405459,939850026,4230372739,1256647217,3585393996,2697235454,3039962178,2933371749,
	// 1562036572,152705284,3961501990,28564985,1101020713,1432483495,3805825671,2896634127,1987054403,3171273214,4027553134,1040338390,2962478315,3768914883,3221538660,3208552736,909582899,138754863,3297363820,
	// 470099579,2024731329,1905310103,1417796907,2335849430,2016672201,783804624,2042375367,3558726676,4047488566,3878875171,4042798972,340938889,1730982572,1932613312,1585165993,1509968046,3055896637,201481366,
	// 1431876315,2431558294,2169233626,2588345649,344557354,2996127788,3630453631,31225136,755899361,1931220684,2613601315,2754642871,1155684340,1039804736,2342987538,1956777345,2647124992,878400658,1752227733,
	// 115936224,132814452,2448435349,4225531806,2265182593,2556267260,3683765494,2381349843,425533424,2105089962,599099141,3316082058,1516687693,579581366,4238700261,74373494,3537953648,1887804568,2369570309,
	// 4161748037,2583156587,3117111058,1997934868,1230719745,3992087881,3256288921,1477262968,3943122032,3099767116,3450669748,2726944001,3235783165,800029767,721323216,3175307064,2778821534,3702275314,1112971016,
	// 1208464815,2886316208,1197514234,1361211748,879533422,1096257750,2210090913,1638823499,537940129,2936408534,2118889604,1156977107,1822089315,432777136,2607169805,391623161,1418375512,2808352364,1538423234,
	// 2462197110,3204164858,2303084321,3247959675,2161672894,131860109,2334275813,98304123,1801330003,2590172193,1288295047,1444361711,3048068295,4111628999,1408863821,2320748645,2799562440,1695510360,472968240,
	// 3382835636,3303158096,2050198131,3724678364,3350142009,3831548606,567207130,3702090536,3346662700,232554249,2849042809,1975348341,3871055944,4221139588,1448316886,204413886,2734139415,697700289,2028850513,
	// 508387185,2300673510,3480763421,1085179555,1622701475,3637349459,4130470606,3215709713,3844510573,41660370,3097311270,2756551288,2339498684,4276691856,3927253891,952553681,3167866960,87110308,696670396,
	// KoZKey =
	// 361314247,1162429539,1977880375,1659336784,2165413738,3400303453,2025084629,2932099871,2172908484,861633009,686846660,3930664248,1558235790,4293777226,2030164768,491126640,569372140,2372877270,169814653,
	// 787816978,3224820645,2119699742,3538562660,3048668038,329814600,2027386413,3100025129,1811354317,1121923984,3081309441,1988237907,2027234121,4104272351,1576762591,1800746683,770898804,1214695468,2559804055,
	// 784832910,2338364994,996043712,3183972336,1612184225,3973278464,510312047,1408382332,1455112834,484252755,4188064208,3808942696,3167582072,3640809019,4122189585,3452612117,1218126785,2585157541,3207804836,
	// 2987432617,2469896172,2683461503,3557988362,2726356220,2578952032,4242665779,3997433007,1791872114,2578206725,1334305774,2533886367,1924168456,2356253966,2171610238,1213332220,3625774784,1414226215,4120600515,
	// 666006537,1713545594,4039875632,3505979526,2358300498,2070163679,1909686440,2901777422,2734607199,1580509836,1747500039,2517588319,4236735283,2102622046,3105497780,3121259132,1369024616,2130156805,261718516,
	// 2382603680,1277925276,349320634,2301008088,1660567391,2162524044,3051719369,4113483893,3465512599,4092889601,763302874,2102638949,2858684788,1273592232,2877080283,2205723915,191823432,3887349872,1621757915,
	// 2989786618,3052889187,1253697053,1177056746,3215813773,2575127333,801553359,1569899497,2364416846,4261794075,1606005605,2688179297,3111449641,1389498449,845851250,2351639132,423532655,3723550833,1422977975,
	// 3661537251,2297239684,976304830,3611188800,3265287165,3581803835,1685880927,4013996934,1028568173,1128569287,3710038763,4192800237,1873976909,1269932918,358198480,522018900,1537209364,3308736866,2554823406,
	// 33200328,1636448833,1848957663,3594812772,1418468300,4080700668,1512588066,1639108583,776524413,505523196,2648460999,1991464848,3859701287,4210150381,1533206818,511392119,3218272472,18666328,2662459861,
	// 1306571374,3582476736,2020646102,2935429318,4055954448,976135523,2463933792,1817651830,3614454258,725380706,3613511057,2762679489,236329806,1826819254,3260800385,255591018,3128946014,2161485584,815701425,
	// 1869915128,3705537516,1424308798,1362332847,3897500207,1765724101,813462993,510181011,389617055,71398813,710973612,808337541,2870098669,4067472102,3516209995,2396650035,1220872527,492791731,3533150342,
	// 2677670209,210909395,3169707910,1265546679,3065797821,4279944890,2537898787,2686445818,3280386720,2732859125,708384979,1419035835,935025562,843050901,3756364032,1935586773,3548406131,2293239556,3272100989,
	// 2582594962,3205915271,1699449978,1219654314,340884135,859404721,2552905507,98686287,4090213177,1551350709,1990117899,3341717825,2192161816,1613668844,4270520760,1019313000,3715354724,3445770873,3687763181,
	// 590631035,2752001943,503366023,2727978621,1290807847,1862240192,356817499,3308062979,405640973,31917169,339179085,3066247744,682030282,2871403976,867228901,4153885015,3347116222,90699976,1405258845,
	// 1113304108,126151259,3544770392,1381916586,2275985685,1932772152,4149075675,4176039982,3092422042,4077631635,493029992,2272538294,718748378,2858492692,3082579820,2012189253,526933076,1225603095,816481173,
	// 483419654,336748150,369136437,1303504898,734568386,2684685564,1579161484,3628666242,3922994716,822433346,2921397537,3710845850,3578576956,1422971419,2999000224,2331236068,1146141718,3367943985,2011746419,
	// 4092207390,3116811404,844413760,403982810,2598087805,3519104580,1942424255,1966328024,1382423311,4103133004,3358429721,859096271,938030143,675113188,330105535,3218366039,2755043850,1314284823,2680896759,
	// 2183128379,1938655615,3287120943,3015145045,3211035021,2764673899,391324404,384577543,836055153,2951597150,660227021,118297137,2222449138,1636648522,3301678523,1355312821,2906247156,3338772981,3581350882,
	// 971006894,4096510989,1970585591,1485807225,1807638416,4108902188,2378124414,4248346977,1167017475,2489370882,25788855,3911310333,795709709,463135905,1001619398,30585747,4250438722,393342929,2216514690,
	// WhiteToPlayZKey = 3234486172
}
