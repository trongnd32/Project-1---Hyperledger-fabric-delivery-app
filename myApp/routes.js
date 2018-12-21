//SPDX-License-Identifier: Apache-2.0

var trans = require('./controller.js');

module.exports = function(app){

  app.get('/get_trans/:id', function(req, res){
    trans.get_trans(req, res);
  });
  app.get('/add_trans/:trans', function(req, res){
    trans.add_trans(req, res);
  });
  app.get('/get_all_trans', function(req, res){
    trans.get_all_trans(req, res);
  });
  app.get('/update_state/:trans', function(req, res){
    trans.update_state(req, res);
  });
}
