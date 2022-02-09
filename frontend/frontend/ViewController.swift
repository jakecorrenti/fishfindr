//
//  ViewController.swift
//  frontend
//
//  Created by Jake Correnti on 12/8/21.
//

import UIKit
import CoreLocation

class ViewController: UIViewController, CLLocationManagerDelegate {
    
    lazy var button: UIButton = {
        let view = UIButton(type: .system)
        view.backgroundColor = .systemGreen
        view.setTitle("Fish Caught!", for: .normal)
        view.setTitleColor(.white, for: .normal)
        view.addTarget(self, action: #selector(buttonPressed), for: .touchUpInside)
        view.layer.cornerRadius = (self.view.bounds.width / 2) / 2
        return view
    }()

    override func viewDidLoad() {
        super.viewDidLoad()
        
        var locationManager: CLLocationManager?
        locationManager = CLLocationManager()
        locationManager?.delegate = self
        locationManager?.requestAlwaysAuthorization()
        
        // Do any additional setup after loading the view.
        button.translatesAutoresizingMaskIntoConstraints = false
        view.addSubview(button)
        NSLayoutConstraint.activate([
            button.widthAnchor.constraint(equalToConstant: view.bounds.width / 2),
            button.heightAnchor.constraint(equalToConstant: view.bounds.width / 2),
            button.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            button.centerYAnchor.constraint(equalTo: view.centerYAnchor)
        ])
    }
    
    @objc
    func buttonPressed() {
        let url = URL(string: "https://671f-151-203-72-221.ngrok.io/api/v1/location")!
        var request = URLRequest(url: url)
        
        let body: [String: Any] = [
            "id" : UUID().uuidString,
                    "latitude": CLLocationManager().location!.coordinate.longitude,
                    "longitude": CLLocationManager().location!.coordinate.latitude,
                    "timestamp": "\(CLLocationManager().location!.timestamp)",
                   ]
        let bodyData = try? JSONSerialization.data(
            withJSONObject: body,
            options: []
        )

        // Change the URLRequest to a POST request
        request.httpMethod = "POST"
        request.httpBody = bodyData
        
        let username = "jcorrenti13"
        let password = "<password>"
        let authData = (username + ":" + password).data(using: .utf8)!.base64EncodedString()
        request.addValue("Basic \(authData)", forHTTPHeaderField: "Authorization")
        
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let session = URLSession.shared
        let task = session.dataTask(with: request) { (data, response, error) in

            if let error = error {
                // Handle HTTP request error
                print(error)
            } else if let data = data {
                // Handle HTTP request response
                print(data)
            } else {
                // Handle unexpected error
            }
        }
        
        task.resume()
    }
}

