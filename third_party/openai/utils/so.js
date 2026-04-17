const J=0,G=1,W=2,z=3,H=4,V=5,B=6,Z=24,K=7,Q=8,Y=9,X=10,tt=11,nt=12,et=13,rt=14,ot=15,it=16,ct=17,st=18,ut=19,at=23,ft=20,lt=21,dt=22,ht=25,pt=26,mt=27,gt=28,wt=29,yt=30,vt=33,bt=34,kt=35,St=new Map;
let Ct = 0, At = Promise.resolve();
{
    // const originalSet = St.set
    // St.set = function(n, e) {
    //     // console.info("St.set called with:", n, e)
    //     // if (n === 50.22 && typeof e === 'object') {
    //     //     e = new Proxy(e, {
    //     //         get(target, prop) {
    //     //             const value = target[prop]
    //     //             console.info("Accessing property:", prop, "value:", value)
    //     //             return value
    //     //         },
    //     //         set(target, prop, value) {
    //     //             console.info("Setting property:", prop, "to value:", value)
    //     //             target[prop] = value
    //     //             return true
    //     //         }
    //     //     })
    //     // }
    //     // if (e < 0) {
    //     //     console.warn("Warning: St.set called with negative value:", n, e)
    //     // }
    //     return originalSet.apply(this, [n, e])
    // }
    // const originalGet = St.get
    // St.get = function(n) {
    //     const value = originalGet.apply(this, [n])
    //     // console.info("St.get called with:", n, "returning:", value)
    //     return value
    // }
}

performance.now = () =>  92357.79999999702
// Date.now = () => 1776434686793
const startTime = Date.now() - 1234 * 35
const oaiData = {
    // 69205.19999998808 + 23152.60000000894
    // 开始的performance.now
  __oai_so_s: 23152.60000000894,
  __oai_so_k: 3,
  __oai_so_kp: 0,
  __oai_so_we: 0,
  __oai_so_wb: 0,
  // 开始的Date.now
  __oai_so_t0: startTime,
  __oai_so_p: 86623.10000000894,
  __oai_so_pc: 3,
  __oai_so_m: 88939.10000000894,
  __oai_so_i: 386,
  __oai_so_ht: 1126.4000000208616,
  __oai_so_hc: 7,
  __oai_so_ss: 162331.38052577656,
  __oai_so_ss2: 242877231.57583714,
  __oai_so_sn: 341,
  __oai_so_cs: 7167.200000077486,
  __oai_so_cs2: 1581126.6199952068,
  __oai_so_cn: 341,
  __oai_so_st: 18,
  __oai_so_sw: 37,
  __oai_so_sp: 7.400527711567218,
  __oai_so_spt: 9,
  __oai_so_sx0: 377.1429138183594,
  __oai_so_sy0: 510.71429443359375,
  __oai_so_lx: 383.4285888671875,
  __oai_so_ly: 514.1429443359375,
  __oai_so_fs: 88530.60000000894,
  __oai_so_fs2: 7505673879.981541,
  __oai_so_fn: 3,
  __oai_so_bc: 0,
  __oai_so_bm: 0,
};
oaiData.__oai_so_h = () => {
    console.info('-------------__oai_so_h called');
}
oaiData.__oai_so_hi = () => {
    console.info('-------------__oai_so_hi called');
}
oaiData.__oai_so_hp = () => {
    console.info('-------------__oai_so_hp called');
}
oaiData.__oai_so_hw = () => {
    console.info('-------------__oai_so_hw called');
}

const createFake = (name, data) => {
  return new Proxy(data, {
    get: (target, prop) => {
      if (name === 'window')
        if (oaiData[prop] !== undefined && oaiData[prop] !== null) {
            target[prop] = oaiData[prop];
            return oaiData[prop]
        } 
      console.info(name, "get:", prop, target[prop]);
      if (target[prop] !== undefined) return target[prop];
      if (name === 'element:div' && prop === 'then') return undefined; // 避免被当做 Promise 处理
      if (name === 'location' && prop === Symbol.toPrimitive) return undefined
      if (name === 'loaderData' && prop === 'root') return undefined
      
      if (name === 'window' && prop?.startsWith?.('__oai_so_')) return undefined
      throw new ReferenceError(name + "." + prop + " is not defined");
    },
    set: (target, prop, value) => {
      console.info(name, "set:", prop, value);
      target[prop] = value;
      return true;
    }
  });
}
const window = createFake("window", {
  Reflect,
  Object,
  Math: createFake("Math", Math),
  Date,
  performance,
  eventMap: {},
  addEventListener: (...args) => {
    console.info("addEventListener called", ...args);
    const eventType = args[0];
    const callback = args[1];
    window.eventMap[eventType] = callback;
  },
  removeEventListener: (...args) => {
    console.info("removeEventListener called", ...args);
  },
})
async function _t(){
    for(;St['get'](Y).length>0;){
        const [n,...e] = St['get'](Y).shift()
        const func = St.get(n)
        // console.info("calling function:", func, "with args:", ...e)
        const r = func(...e);
        // if (typeof r === 'number' || typeof r === 'string' || typeof r === 'boolean')
        //     console.info("function result:",r)
        r&&typeof r['then']==='function'&&await r,Ct++
    }
}
function Rt(t,n){
    let r="";
    for(let o=0;o<t['length'];o++)
        r+=String['fromCharCode'](t['charCodeAt'](o)^n['charCodeAt'](o%n['length']));
    
    console.info('decode result:', [t, n, r])
    return r
}
function jt(t,n){
    return new Promise(((e,r) => {
    void 0 !== n &&
        ((function () {
        // St = new Map
        (St["clear"](),
            St.set(J, Nt),
            St["set"](G, (n, e) =>
                St["set"](n, Rt("" + St["get"](n), "" + St.get(e))),
            ),
            St["set"](W, (n, e) => St["set"](n, e)),
            St["set"](V, (n, e) => {
            const o = St['get'](n);
            Array['isArray'](o)
                ? o.push(St['get'](e))
                : St.set(n, o + St.get(e));
            }),
            St.set(mt, (n, e) => {
            const o = St.get(n);
            Array.isArray(o)
                ? o["splice"](o["indexOf"](St.get(e)), 1)
                : St["set"](n, o - St["get"](e));
            }),
            St.set(wt, (n, e, r) => St.set(n, St["get"](e) < St["get"](r))),
            St.set(vt, (n, e, r) => {
            const i = Number(St["get"](e)),
                c = Number(St["get"](r));
            St["set"](n, i * c);
            }),
            St.set(kt, (n, e, r) => {
            const i = Number(St.get(e)),
                c = Number(St["get"](r));
            St["set"](n, 0 === c ? 0 : i / c);
            }),
            St["set"](B, (n, e, r) => St.set(n, St["get"](e)[St.get(r)])),
            St.set(K, (n, ...e) =>{
                const func = St["get"](n)
                const args = e["map"]((n) => St["get"](n))
                // console.info("calling function:", func, "with args:", ...args)
                const result = func(...args)
                // console.info("function result:", result)
                return result
            }
            ),
            St["set"](ct, (n, e, ...r) => {
            try {
                const args = r["map"]((t) => St["get"](t))
                const func = St.get(e)
                // console.info("calling function:", func, "with args:", ...args)
                const t = func(...args);
                if (t && typeof t.then === "function")
                    return t["then"]((t) => {
                        St.set(n, t);
                    }).catch((t) => {
                        St.set(n, "" + t);
                    });
                // console.info("function result:", t)
                St["set"](n, t);
            } catch (t) {
                St["set"](n, "" + t);
            }
            }),
            St["set"](et, (n, e, ...r) => {
            try {
                const func = St["get"](e)
                console.info("calling function:", func, "with args:", ...r)
                func(...r);
            } catch (t) {
                St["set"](n, "" + t);
            }
            }),
            St.set(Q, (n, e) => St["set"](n, St["get"](e))),
            St["set"](X, window),
            St["set"](tt, (n, e) =>
            St["set"](
                n,
                (Array.from(document.scripts || [])
                ["map"]((n) => n?.["src"]?.["match"](St["get"](e)))
                ["filter"]((n) => n?.["length"])[0] ?? [])[0] ?? null,
            ),
            ),
            St["set"](nt, (n) => St["set"](n, St)),
            St["set"](rt, (n, e) =>
            St["set"](n, JSON["parse"]("" + St["get"](e))),
            ),
            St.set(ot, (n, e) => St["set"](n, JSON.stringify(St.get(e)))),
            St["set"](st, (n) => St["set"](n, atob("" + St["get"](n)))),
            St["set"](ut, (n) => St["set"](n, btoa("" + St["get"](n)))),
            St["set"](ft, (n, e, r, ...o) =>{
                if (St["get"](n) === St["get"](e)) {
                    const func = St["get"](r)
                    console.info("condition met, calling function:", func, "with args:", ...o)
                    const result = func(...o);
                    // console.info("function result:", result)
                    return result
                }
                return null
            }
            ),
            St["set"](lt, (n, e, r, o, ...i) =>
            Math["abs"](St["get"](n) - St.get(e)) > St["get"](r)
                ? St["get"](o)(...i)
                : null,
            ),
            St["set"](at, (n, e, ...r) =>
            void 0 !== St["get"](n) ? St["get"](e)(...r) : null,
            ),
            St["set"](Z, (n, e, r) =>
            St["set"](n, St.get(e)[St.get(r)]["bind"](St.get(e))),
            ),
            St["set"](bt, (n, e) => {
            try {
                const t = St["get"](e);
                return Promise["resolve"](t)["then"]((t) => {
                St.set(n, t);
                });
            } catch {
                return;
            }
            }),
            St["set"](dt, (n, e) => {
            const o = [...St["get"](Y)];
            return (
                St["set"](Y, [...e]),
                _t()
                ["catch"]((t) => {
                    St["set"](n, "" + t);
                })
                ["finally"](() => {
                    St["set"](Y, o);
                })
            );
            }),
            St["set"](gt, () => {}),
            St.set(pt, () => {}),
            St["set"](ht, () => {}));
        })(),
        (Ct = 0),
        St.set(it, n));
        let i=!1;
        const c=setTimeout((()=>{!i&&(i=!0,r(new Error('session_observer_vm_timeout')))}),6e4),s=t=>{i||(i=!0,clearTimeout(c),e(t))};St['set'](z,(t=>{s(btoa(""+t))})),St.set(H,(t=>{(t=>{i||(i=!0,clearTimeout(c),r(t))})(btoa(""+t))})),
        St['set'](yt,((t,n,e,r)=>{
            const c=Array['isArray'](r),s=c?e:[],u=(c?r:e)||[];
            St['set'](t,((...t)=>{
                const r=[...St['get'](Y)];
                if(c)
                    for(let n=0;n<s['length'];n++){
                        const r=s[n],o=t[n];
                        St['set'](r,o)
                    }
                return St['set'](Y,[...u]),
                _t()['then']((()=>St.get(n)))
                ['catch']((t=>""+t))
                ['finally']((()=>{
                    St.set(Y,r)
                }))
            }))
        }));
        try{
        St['set'](Y,JSON['parse'](Rt(atob(t),""+St['get'](it)))),_t()['catch']((t=>{s(btoa(Ct+": "+t))}))
        }catch(t){
        s(btoa(Ct+": "+t))
        }
    }))
}
function Ot(t){
    const e=At['then'](t,t);
    return At=e['then']((()=>{}),(()=>{})),e
}
function Nt(snapshot_dx){
    return Ot((()=>jt(snapshot_dx)))
};
function Et(collector_dx, p){
    return Ot((()=>jt(collector_dx, p)))
};
module.exports = {
    /**
     * 
     * @param {string} p https://sentinel.openai.com/backend-api/sentinel/req 请求体p字段
     * @param {string} collector_dx backend-api/sentinel/req 响应体中的字段
     * @param {string} snapshot_dx backend-api/sentinel/req 响应体中的字段
     * @returns 
     */
    so: async (p, collector_dx, snapshot_dx) => {
        await Et(collector_dx, p)
        return await Nt(snapshot_dx)
    }
}
// (async () => {
//     // 请求体的p字段
//     const p = "gAAAAACWzI4NTQsIkZyaSBBcHIgMTcgMjAyNiAyMjowMzozMiBHTVQrMDgwMCAo5Lit5Zu95qCH5YeG5pe26Ze0KSIsNDI5NDk2NzI5NiwxNywiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzE0Ni4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiaHR0cHM6Ly9zZW50aW5lbC5vcGVuYWkuY29tL2JhY2tlbmQtYXBpL3NlbnRpbmVsL3Nkay5qcyIsbnVsbCwiemgtQ04iLCJ6aC1DTix6aCIsNCwibWFuYWdlZOKIkltvYmplY3QgTmF2aWdhdG9yTWFuYWdlZERhdGFdIiwibG9jYXRpb24iLCJvbnBvaW50ZXJlbnRlciIsMTc4MTkuNzk5OTk5OTk3MDIsIjAyMWI4NDQ2LWI2YzUtNDIyOC05ZmNmLTgxOTIwNGFiNDlmZSIsIiIsMjIsMTc3NjQzNDU5NDQzNi4yLDAsMCwwLDAsMCwwLDBd~S"
//     const collector_dx = "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"
//     const snapshot_dx = "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"
//     await Et(collector_dx, p)
//     // for (const k in oaiData) {
//     //     window[k] = oaiData[k];
//     // }
//     const so = await Nt(snapshot_dx)
//     console.log('so:', so)
//     console.info(so === 'QxobHAcBGhQQdV99TXNWaW9tYFJ0e293c35PQXhyeXkTEB0aChwcAgAaFBByTwUTEB0aDhwcAAkaFBByeQUTEB0aARwcAgkaFBByeQUTEB0aDxccBwoaFBByamFKcHBTdmpwaHx5b3N4fl9ndnJpTxMQHRoMHxwGDxoUEHJ5BRMQHRoLGBwCDBoUEHJ5BRMQHRoAAAoIGgIMcXAFBQweEwAIAAAHGgIMcWAAYmxzZHR8bFV1el8ec3dfBRMQHRoKHxwFCxoUEHJpAHRwcG10anBWc3pZCnNwaRMPExQaHwMfDgkMCBN5eRprcHlxeXFjDHtvY1ZxcEl7e3lpWQ8TFBoYAR8ACgwIE3tPEw8TFBofBB8JDQwIE3tPb2sTFBobAx8JDgwIE3tpRVFzegBjcXMMe29jVnFwb3t7fF9FaxMUGhsGHw4ODAgTfE8TDxMUGhYBHwsNDAgTe2p3U3N5cW92SWF+b0VwcX9/e3l7eV1TExQaFwIfDA0MCBN7QHdXcl95YHZjDHNvVUJ1f39Bd3tfRVATFBoYAh8MAQwIE3t6d1NyX1tjc3Nff21jBXV+SWN3eXkaUBNF')
// })();

/*
window.tempData = {}
const keys = ['__oai_so_h', '__oai_so_hi', '__oai_so_hp', '__oai_so_hw', '__oai_so_s', '__oai_so_k', '__oai_so_kp', '__oai_so_we', '__oai_so_wb', '__oai_so_wl', '__oai_so_t0', '__oai_so_p', '__oai_so_pc', '__oai_so_m', '__oai_so_i', '__oai_so_ht', '__oai_so_hc', '__oai_so_ss', '__oai_so_ss2', '__oai_so_sn', '__oai_so_cs', '__oai_so_cs2', '__oai_so_cn', '__oai_so_st', '__oai_so_sw', '__oai_so_sp', '__oai_so_spt', '__oai_so_sx0', '__oai_so_sy0', '__oai_so_lx', '__oai_so_ly', '__oai_so_fs', '__oai_so_fs2', '__oai_so_fn', '__oai_so_bc', '__oai_so_bm']
const cfg = Object.fromEntries(keys.map(e => ([e, {
    get() { 
        const result = window.tempData[e] 
        console.log(`get ${e}:`, result)
        return result
    },
    set(v) {
        if (v === null)return
        window.tempData[e] = v }
}])))
keys.forEach(e => {
    window.tempData[e] = window[e]
})
Object.defineProperties(window, cfg)
*/