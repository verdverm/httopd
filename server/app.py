#!flask/bin/python
from flask import Flask, make_response
import os

app = Flask(__name__)


@app.errorhandler(404)
def not_found(error):
	return make_response('Nada Dawg!' , 404)

@app.errorhandler(400)
def not_found(error):
	print error
	return make_response('What you say Dawg?!', 400)

@app.route('/')
@app.route('/index')
def index():
    return "Yo Dawg!"

@app.route('/page1')
@app.route('/page1/<string:subpage>')
@app.route('/page1/<string:subpage>/<string:subsubpage>')
def page1(subpage = "", subsubpage = ""):
    res = "Yo Dawg! You found page1"
    if (subpage is not "" ):
        res += "/" + subpage
    if (subsubpage is not "" ):
    	res += "/" + subsubpage
    return res

@app.route('/page2')
@app.route('/page2/<string:subpage>')
@app.route('/page2/<string:subpage>/<string:subsubpage>')
def page2(subpage = "", subsubpage = ""):
    res = "Yo Dawg! You found page2"
    if (subpage is not "" ):
        res += "/" + subpage
    if (subsubpage is not "" ):
    	res += "/" + subsubpage
    return res

@app.route('/page3')
@app.route('/page3/<string:subpage>')
@app.route('/page3/<string:subpage>/<string:subsubpage>')
def page3(subpage = "", subsubpage = ""):
    res = "Yo Dawg! You found page3"
    if (subpage is not "" ):
        res += "/" + subpage
    if (subsubpage is not "" ):
    	res += "/" + subsubpage
    return res

@app.route('/page4')
@app.route('/page4/<string:subpage>')
@app.route('/page4/<string:subpage>/<string:subsubpage>')
def page4(subpage = "", subsubpage = ""):
    res = "Yo Dawg! You found page4"
    if (subpage is not "" ):
        res += "/" + subpage
    if (subsubpage is not "" ):
    	res += "/" + subsubpage
    return res

name = ""

if __name__ == '__main__':
    port = int(os.environ['PORT'])
    app.run(host='0.0.0.0', port=port, debug=True)

