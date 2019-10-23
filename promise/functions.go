package promise

func New(resolver func(resolve, reject func(interface{}))) *Promise {
	p := &Promise{pending, nil, make([]*Handler, 0)}
	go doResolve(resolver, p.resolve, p.reject)
	return p
}

func doResolve(fn func(_, _ func(value interface{})), onFulfilled, onRejected func(value interface{})) {
	done := false
	defer func() {
		err := recover()
		if err != nil {
			if done {
				return
			}
			done = true
			onRejected(err)
		}
	}()

	fn(func(value interface{}) {
		if done {
			return
		}

		done = true
		onFulfilled(value)
	}, func(reason interface{}) {
		if done {
			return
		}

		done = true
		onRejected(reason)
	})
}

func getThen(value interface{}) (func(onFulfilled, onRejected func(reason interface{})), bool) {
	promise, ok := value.(Promise)
	if ok {
		return func(onFulfilled, onRejected func(reason interface{})) {
			resolve := func(value interface{}) interface{} {
				if onFulfilled != nil {
					onFulfilled(value)
				}
				return nil
			}

			reject := func(value interface{}) interface{} {
				if onRejected != nil {
					onRejected(value)
				}
				return nil
			}
			promise.thenInternal(resolve, reject)
		}, true
	}

	return nil, false
}

func isPromise(value interface{}) bool {
	_, ok := value.(*Promise)
	return ok
}

func makeFulfillChain(result interface{}, onFulfilled func(value interface{}) interface{}, onRejected func(reason interface{}) interface{}, resolve func(interface{}), reject func(interface{})) {
	if isPromise(result) {
		result.(*Promise).done(func(value interface{}) {
			makeFulfillChain(value, onFulfilled, onRejected, resolve, reject)
		}, func(reason interface{}) {
			makeRejectChain(reason, onFulfilled, onRejected, resolve, reject)
		})
		return
	}

	if onFulfilled != nil {
		makeFulfillChain(onFulfilled(result), nil, nil, resolve, reject)
		return
	}

	resolve(result)
}

func makeRejectChain(reason interface{}, onFulfilled func(value interface{}) interface{}, onRejected func(reason interface{}) interface{}, resolve func(interface{}), reject func(interface{})) {
	if isPromise(reason) {
		reason.(*Promise).done(func(value interface{}) {
			makeFulfillChain(reason, onFulfilled, onRejected, resolve, reject)
		}, func(reason interface{}) {
			makeRejectChain(reason, onFulfilled, onRejected, resolve, reject)
		})
		return
	}

	if onRejected != nil {
		makeFulfillChain(onRejected(reason), nil, nil, resolve, reject)
		return
	}

	reject(reason)
}
