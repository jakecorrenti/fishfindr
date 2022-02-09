//
// Created by Jake Correnti on 1/22/22.
//

import UIKit

class TabController: UITabBarController {

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = .systemBackground
        let button = UINavigationController(rootViewController: ViewController())
        button.tabBarItem = UITabBarItem(title: "Button", image: UIImage.init(systemName: "archivebox", withConfiguration: UIImage.SymbolConfiguration(weight: .semibold)), tag: 0)

        let map = UINavigationController(rootViewController: MapVC())
        button.tabBarItem = UITabBarItem(title: "Map", image: UIImage.init(systemName: "mappin.and.ellipse", withConfiguration: UIImage.SymbolConfiguration(weight: .semibold)), tag: 1)

        viewControllers = [map, button]

        tabBar.tintColor = .systemGreen

    }

}